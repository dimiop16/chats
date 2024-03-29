import { checkResponse, fetchWithTimeout } from './fetchUtil';

export async function getChannelTopic(server, token, channelId, topicId) {
  let topic = await fetchWithTimeout(`https://${server}/content/channels/${channelId}/topics/${topicId}/detail?agent=${token}`, 
    { method: 'GET' });
  checkResponse(topic)
  return await topic.json()
}

