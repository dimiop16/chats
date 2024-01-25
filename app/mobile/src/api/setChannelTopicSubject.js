import { checkResponse, fetchWithTimeout } from './fetchUtil';

export async function setChannelTopicSubject(server, token, channelId, topicId, dataType, data) {
  let subject = { data: JSON.stringify(data, (key, value) => {
    if (value !== null) return value
  }), datatype: dataType };

  let channel = await fetchWithTimeout(`https://${server}/content/channels/${channelId}/topics/${topicId}/subject?agent=${token}&confirm=true`,
    { method: 'PUT', body: JSON.stringify(subject) });
  checkResponse(channel);
}
