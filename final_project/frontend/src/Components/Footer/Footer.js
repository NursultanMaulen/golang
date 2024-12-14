import { Layout, Row, Col, Typography, Button } from "antd";
import { FaGithub, FaLinkedin } from "react-icons/fa";
import { Link } from "react-router-dom";

const { Footer } = Layout;
const { Text } = Typography;

function CustomFooter() {
  return (
    <Footer
      style={{
        backgroundColor: "#f0f2f5",
        padding: "20px 50px",
        position: "fixed",
        bottom: 0,
        left: 0,
        right: 0,
        zIndex: 1001,
      }}
    >
      <Row justify="space-between" align="middle">
        <Col>
          <Button type="link" style={{ padding: 5 }}>
            <Link to="/">Home</Link>
          </Button>
          <Button type="link" style={{ padding: 5 }}>
            <Link to="/explore">Explore</Link>
          </Button>
          <Button type="link" style={{ padding: 5 }}>
            <Link to="/about">About Us</Link>
          </Button>
        </Col>

        <Col>
          <Button
            type="link"
            href="https://github.com/NursultanMaulen"
            target="_blank"
            icon={<FaGithub style={{ fontSize: "24px" }} />}
          />

          <Button
            type="link"
            href="https://www.linkedin.com/in/nmaulen/"
            target="_blank"
            icon={
              <FaLinkedin style={{ fontSize: "24px", marginLeft: "20px" }} />
            }
          />
        </Col>

        <Col>
          <Text style={{ textAlign: "center" }}>
            © {new Date().getFullYear()} Maulen Nursultan. All rights reserved.
          </Text>
        </Col>
      </Row>
    </Footer>
  );
}

export default CustomFooter;
