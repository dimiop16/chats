import { useState, useRef } from 'react';
import axios from 'axios';
import Resizer from "react-image-file-resizer";

const ENCRYPTED_BLOCK_SIZE = (256 * 1024);
const IMAGE_SCALE_SIZE = (128 * 1024);
const GIF_TYPE = 'image/gif';
const WEBP_TYPE = 'image/webp';

export function useUploadContext() {

  const [state, setState] = useState({
    progress: new Map(),
  });
  const channels = useRef(new Map());
  const index = useRef(0);

  const updateState = (value) => {
    setState((s) => ({ ...s, ...value }));
  };

  const updateComplete = (channel, topic) => {
    let topics = channels.current.get(channel);
    if (topics) {
      topics.delete(topic);
    }
    updateProgress();
  }

  const updateProgress = () => {
    let progress = new Map();
    channels.current.forEach((topics, channel) => {
      let assets = [];
      topics.forEach((entry, topic, map) => {
        let active = entry.active ? 1 : 0;
        assets.push({
          upload: entry.index,
          topicId: topic,
          active: entry.active,
          uploaded: entry.assets.length,
          index: entry.assets.length + active,
          count: entry.assets.length + entry.files.length + active,
          error: entry.error,
        });
      });
      if (assets.length) {
        progress.set(channel, assets.sort((a, b) => (a.upload < b.upload) ? 1 : -1));
      }
    });
    updateState({ progress });
  }

  const abort = (channelId, topicId) => {
    const channel = channels.current.get(channelId);
    if (channel) {
      const topic = channel.get(topicId);
      if (topic) {
        topic.cancel.abort();
        channel.delete(topicId);
        updateProgress();
      }
    }
  }

  const actions = {
    addTopic: (token, channelId, topicId, files, success, failure, contact) => {
      if (contact) {
        const { server, cardId } = contact;

        let host = "";
        if (server) {
          host = `https://${server}`
        }

        const controller = new AbortController();
        const entry = {
          index: index.current,
          baseUrl: `${host}/content/channels/${channelId}/topics/${topicId}/`,
          urlParams: `?contact=${token}`,
          files,
          assets: [],
          current: null,
          error: false,
          success,
          failure,
          cancel: controller,
        }
        index.current += 1;
        const key = `${cardId}:${channelId}`;
        if (!channels.current.has(key)) {
          channels.current.set(key, new Map());
        }
        const topics = channels.current.get(key);
        topics.set(topicId, entry);
        upload(entry, updateProgress, () => { updateComplete(key, topicId) });
      }
      else {
        const controller = new AbortController();
        const entry = {
          index: index.current,
          baseUrl: `/content/channels/${channelId}/topics/${topicId}/`,
          urlParams: `?agent=${token}`,
          files,
          assets: [],
          current: null,
          error: false,
          success,
          failure,
          cancel: controller,
        }
        index.current += 1;
        const key = `:${channelId}`;
        if (!channels.current.has(key)) {
          channels.current.set(key, new Map());
        }
        const topics = channels.current.get(key);
        topics.set(topicId, entry);
        upload(entry, updateProgress, () => { updateComplete(key, topicId) } );
      }
    },
    cancelTopic: (channelId, topicId, cardId) => {
      if (cardId) {
        abort(`${cardId}:${channelId}`, topicId);
      }
      else {
        abort(`:${channelId}`, topicId);
      }
    },
    clearErrors: (cardId, channelId) => {
      const key = cardId ? `${cardId}:${channelId}` : `:${channelId}`;
      const topics = channels.current.get(key);
      if (topics) {
        topics.forEach((topic, topicId) => {
          if (topic.error) {
            topic.cancel.abort();
            topics.delete(topicId);
            updateProgress();
          }
        });
      }
    },
    clear: () => {
      channels.current.forEach((topics, channelId) => {
        topics.forEach((assets, topicId) => {
          assets.cancel.abort();
        });
      });
      channels.current.clear();
      updateProgress();
    }
  }

  return { state, actions }
}

function getImageThumb(data) {
  return new Promise((resolve, reject) => {
    if ((data.type === GIF_TYPE || data.type === WEBP_TYPE) && data.size < IMAGE_SCALE_SIZE) {
      const reader = new FileReader();
      reader.readAsDataURL(data);
      reader.onload = function () {
        resolve(reader.result);
      };
      reader.onerror = function (error) {
        reject();
      };
    }
    else {
      Resizer.imageFileResizer(data, 192, 192, 'JPEG', 50, 0,
      uri => {
        resolve(uri);
      }, 'base64', 128, 128 );
    }
  });
}

function getVideoThumb(data, pos) {
  return new Promise((resolve, reject) => {
    const url = URL.createObjectURL(data);
    var video = document.createElement("video");
    var timeupdate = function (ev) {
      video.removeEventListener("timeupdate", timeupdate);
      video.pause();
      setTimeout(() => {
        var canvas = document.createElement("canvas");
        if (video.videoWidth > video.videoHeight) {
          canvas.width = 192;
          canvas.height = Math.floor((192 * video.videoHeight / video.videoWidth));
        }
        else {
          canvas.height  = 192;
          canvas.width = Math.floor((192 * video.videoWidth / video.videoHeight));
        }
        canvas.getContext("2d").drawImage(video, 0, 0, canvas.width, canvas.height);
        var image = canvas.toDataURL("image/jpeg", 0.75);
        resolve(image);
        canvas.remove();
        video.remove();
        URL.revokeObjectURL(url);
      }, 1000);
    };
    video.addEventListener("timeupdate", timeupdate);
    video.preload = "metadata";
    video.src = url;
    video.muted = true;
    video.playsInline = true;
    video.currentTime = pos;
    video.play();
  });
}

async function getThumb(data, type, position) {
  
  if (type === 'image') {
    return await getImageThumb(data);
  }
  else if (type === 'video') {
    return await getVideoThumb(data, position);
  }
  else {
    return null;
  }
}

async function upload(entry, update, complete) {
  if (!entry.files?.length) {
    entry.success(entry.assets);
    complete();
  }
  else {
    const file = entry.files.shift();
    entry.active = {};
    try {
      if (file.encrypted) {
        const { size, getEncryptedBlock, position, label, extension, image, video, audio, binary } = file;
        const { data, type } = image ? { data: image, type: 'image' } : video ? { data: video, type: 'video' } : audio ? { data: audio, type: 'audio' } : { data: binary, type: 'binary' }
        const thumb = await getThumb(data, type, position);
        const parts = [];
        for (let pos = 0; pos < size; pos += ENCRYPTED_BLOCK_SIZE) {
          const len = pos + ENCRYPTED_BLOCK_SIZE > size ? size - pos : ENCRYPTED_BLOCK_SIZE;
          const { blockEncrypted, blockIv } = await getEncryptedBlock(pos, len);
          const part = await axios.post(`${entry.baseUrl}blocks${entry.urlParams}`, blockEncrypted, {
            headers: {'Content-Type': 'text/plain'},
            signal: entry.cancel.signal,
            onUploadProgress: (ev) => {
              const { loaded, total } = ev;
              const partLoaded = pos + Math.floor(len * loaded / total);
              entry.active = { loaded: partLoaded, total: size }
              update();
            }
          });
          parts.push({ blockIv, partId: part.data.assetId });
        }
        entry.assets.push({
          encrypted: { type, thumb, label, extension, parts }
        });
      }
      else if (file.image) {
        const formData = new FormData();
        formData.append('asset', file.image);
        let transform = encodeURIComponent(JSON.stringify(["ithumb;photo", "ilg;photo"]));
        let asset = await axios.post(`${entry.baseUrl}assets${entry.urlParams}&transforms=${transform}`, formData, {
          signal: entry.cancel.signal,
          onUploadProgress: (ev) => {
            const { loaded, total } = ev;
            entry.active = { loaded, total }
            update();
          },
        });
        entry.assets.push({
          image: {
            thumb: asset.data.find(item => item.transform === 'ithumb;photo').assetId,
            full: asset.data.find(item => item.transform === 'ilg;photo').assetId, 
          }
        });
      }
      else if (file.video) {
        const formData = new FormData();
        formData.append('asset', file.video);
        let thumb = 'vthumb;video;' + file.position;
        let transform = encodeURIComponent(JSON.stringify(["vlq;video", "vhd;video", thumb]));
        let asset = await axios.post(`${entry.baseUrl}assets${entry.urlParams}&transforms=${transform}`, formData, {
          signal: entry.cancel.signal,
          onUploadProgress: (ev) => {
            const { loaded, total } = ev;
            entry.active = { loaded, total }
            update();
          },
        });
        entry.assets.push({
          video: {
            thumb: asset.data.find(item => item.transform === thumb).assetId,
            lq: asset.data.find(item => item.transform === 'vlq;video').assetId,
            hd: asset.data.find(item => item.transform === 'vhd;video').assetId,
          }
        });
      }
      else if (file.audio) {
        const formData = new FormData();
        formData.append('asset', file.audio);
        let transform = encodeURIComponent(JSON.stringify(["acopy;audio"]));
        let asset = await axios.post(`${entry.baseUrl}assets${entry.urlParams}&transforms=${transform}`, formData, {
          signal: entry.cancel.signal,
          onUploadProgress: (ev) => {
            const { loaded, total } = ev;
            entry.active = { loaded, total }
            update();
          },
        });
        entry.assets.push({
          audio: {
            label: file.label,
            full: asset.data.find(item => item.transform === 'acopy;audio').assetId,
          }
        });
      }
      else if (file.binary) {
        const formData = new FormData();
        formData.append('asset', file.binary);
        let asset = await axios.post(`${entry.baseUrl}blocks${entry.urlParams}&body=multipart`, formData, {
          signal: entry.cancel.signal,
          onUploadProgress: (ev) => {
            const { loaded, total } = ev;
            entry.active = { loaded, total }
            update();
          },
        });
        entry.assets.push({
          binary: {
            label: file.label,
            extension: file.extension,
            data: asset.data.assetId,
          }
        });
      }
      entry.active = null;
      upload(entry, update, complete);
    }
    catch (err) {
      console.log(err);
      entry.failure();
      entry.error = true;
      update();
    }
  }
}

