import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Form, Input, Button, message, Card } from "antd";
import { getVideoById, updateVideo } from "../../services/VideoService";
import { useExplorePageContext } from "../../Context/ExplorepageContext";

function Videopage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { fetchVideos } = useExplorePageContext();
  const [video, setVideo] = useState(null);

  useEffect(() => {
    async function fetchVideo() {
      try {
        const response = await getVideoById(id);
        setVideo(response.data);
      } catch (error) {
        message.error("Failed to fetch video. Please try again.");
      }
    }
    fetchVideo();
  }, [id]);

  const onFinish = async (values) => {
    try {
      await updateVideo(id, values);
      message.success("Video updated successfully!");
      navigate("/explore");
    } catch (error) {
      message.error("Failed to update video. Please try again.");
    }
  };

  if (!video) {
    return <div>Loading...</div>;
  }

  return (
    <Card style={{ maxWidth: "600px", margin: "auto", marginTop: "20px" }}>
      <div
        style={{
          position: "relative",
          paddingBottom: "56.25%",
          height: 0,
          marginBottom: "20px",
        }}
      >
        <iframe
          src={video.videoUrl}
          title={video.title}
          style={{
            position: "absolute",
            top: 0,
            left: 0,
            width: "100%",
            height: "100%",
            border: "none",
            borderRadius: "8px",
          }}
          allowFullScreen
        />
      </div>
      <Form layout="vertical" initialValues={video} onFinish={onFinish}>
        <Form.Item
          label="Title"
          name="title"
          rules={[{ required: true, message: "Please enter the video title!" }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Description"
          name="description"
          rules={[{ required: true, message: "Please enter the description!" }]}
        >
          <Input.TextArea rows={4} />
        </Form.Item>
        <Form.Item
          label="Category"
          name="category"
          rules={[{ required: true, message: "Please select a category!" }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Video URL"
          name="videoUrl"
          rules={[{ required: true, message: "Please enter the video URL!" }]}
        >
          <Input />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit">
            Save Changes
          </Button>
        </Form.Item>
      </Form>
    </Card>
  );
}

export default Videopage;
