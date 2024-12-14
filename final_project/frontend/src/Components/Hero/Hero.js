import React from "react";
import { Link } from "react-router-dom";
import { Layout, Typography, Button, Row, Col, Image } from "antd";
import learningIllustration from "../../assets/undraw_online_learning_re_qw08.svg";

const { Title, Paragraph } = Typography;

function Hero() {
  return (
    <Layout
      style={{
        background: "#fff",
        padding: "64px 24px",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        minHeight: "50vh",
      }}
    >
      <Row justify="space-between" align="middle">
        <Col xs={24} md={12} style={{ textAlign: "center", md: "left" }}>
          <Title level={1} style={{ fontWeight: "bold" }}>
            Discover StreamBox
          </Title>
          <Paragraph style={{ fontSize: "18px" }}>
            Explore new content!
          </Paragraph>
          <Link to="/explore">
            <Button type="primary" size="large">
              Explore Now
            </Button>
          </Link>
        </Col>

        <Col xs={24} md={12} style={{ textAlign: "center" }}>
          <Image
            src={learningIllustration}
            alt="Learning illustration"
            preview={false}
            width="100%"
            style={{ maxWidth: "400px", borderRadius: "8px" }}
          />
        </Col>
      </Row>
    </Layout>
  );
}

export default Hero;
