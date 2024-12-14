import React, { useState } from "react";
import { Button, Input, Form, Typography, Layout, message } from "antd";
import { Link, useNavigate } from "react-router-dom";
import { useLoginSignupContext } from "../../Context/LoginSignupContext";

const { Title, Text } = Typography;
const { Content } = Layout;

function LoginInputs() {
  const { loginUser } = useLoginSignupContext();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [error, setError] = useState("");

  const submitLoginData = async (values) => {
    const { email, password } = values;

    const result = await loginUser(email, password);

    if (result.success) {
      message.success("Login successful!");
      navigate("/explore");
    } else {
      message.error(result.message || "Login failed. Please try again.");
    }
  };

  return (
    <Layout
      style={{
        minHeight: "85vh",
        display: "flex",
        alignItems: "center",
        background: "#fff",
        justifyContent: "center",
      }}
    >
      <Content
        style={{
          maxWidth: "400px",
          width: "100%",
          padding: "20px",
          borderRadius: "8px",
        }}
      >
        <Title level={3} style={{ textAlign: "center" }}>
          Login Page
        </Title>
        <Form layout="vertical" form={form} onFinish={submitLoginData}>
          <Form.Item
            label="Email"
            name="email"
            rules={[{ required: true, message: "Please input your email!" }]}
          >
            <Input
              style={{ borderRadius: "4px" }}
              type="email"
              placeholder="Enter your email"
            />
          </Form.Item>

          <Form.Item
            label="Password"
            name="password"
            rules={[
              {
                required: true,
                message: "Please input your password!",
                min: 6,
              },
            ]}
          >
            <Input.Password
              style={{ borderRadius: "4px" }}
              placeholder="Enter your password"
            />
          </Form.Item>

          {error && <Text type="danger">{error}</Text>}

          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              Login
            </Button>
          </Form.Item>
        </Form>

        <Text style={{ textAlign: "center" }}>
          Not a member?{" "}
          <Link to="/signup" style={{ color: "#1890ff" }}>
            Signup
          </Link>{" "}
          here
        </Text>
      </Content>
    </Layout>
  );
}

export default LoginInputs;
