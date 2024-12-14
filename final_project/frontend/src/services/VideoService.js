import axios from "axios";

const API_URL = "http://localhost:3001/videos";

export const getVideoById = async (id) => {
  return await axios.get(`${API_URL}/${id}`);
};

export const updateVideo = async (id, videoData) => {
  return await axios.put(`${API_URL}/${id}`, videoData);
};
