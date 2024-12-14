import React, { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import { Layout, Typography, Row, Col, Card } from "antd";
import { Footer, Header, Sidebar } from "../../Components/IndexAllComponents";
import { useExplorePageContext } from "../../Context/IndexAllContext";

const { Content } = Layout;
const { Title } = Typography;

function SearchResultsPage() {
  const { state } = useExplorePageContext();
  const { videosdata } = state;
  const location = useLocation();
  const [filteredVideos, setFilteredVideos] = useState([]);

  const searchParams = new URLSearchParams(location.search);
  const query = searchParams.get("query");

  useEffect(() => {
    if (query) {
      const results = videosdata.filter((video) =>
        video.title.toLowerCase().includes(query.toLowerCase())
      );
      setFilteredVideos(results);
    }
  }, [query, videosdata]);

  return (
    <Layout style={{ minHeight: "100vh", overflow: "hidden" }}>
      <Header />
      <Layout>
        <Sidebar />
        <Content
          style={{
            backgroundColor: "#fff",
          }}
        >
          <Title
            style={{
              position: "relative",
              justifyContent: "center",
            }}
            level={2}
          >
            {filteredVideos.length > 0
              ? `Search Results for "${query}"`
              : "No results found."}
          </Title>
          <Row gutter={[16, 16]} justify="center">
            {filteredVideos.map((video) => (
              <Col key={video.id} xs={24} sm={12} md={8} lg={6}>
                <Card
                  title={video.title}
                  hoverable
                  style={{
                    width: "100%",
                    borderRadius: "8px",
                    overflow: "hidden",
                  }}
                >
                  <iframe
                    width="100%"
                    height="200"
                    src={video.videoUrl}
                    frameBorder="0"
                    allowFullScreen
                    title={video.title}
                    style={{ borderRadius: "8px 8px 0 0" }}
                  ></iframe>
                </Card>
              </Col>
            ))}
          </Row>
        </Content>
      </Layout>
      <Footer />
    </Layout>
  );
}

export default SearchResultsPage;
