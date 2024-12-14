import React, { useState } from "react";
import { Button, Form, Input, Typography, Layout, message } from "antd";
import { Link, useNavigate } from "react-router-dom";

const { Title, Text } = Typography;
const { Content } = Layout;

function SignupInputs() {
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [error, setError] = useState("");

  const submitSignupData = async (values) => {
    const userData = {
      username: values.name,
      email: values.email,
      password_hash: values.password,
      role: "user",
    };

    try {
      const response = await fetch("http://localhost:8080/api/users", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(userData),
      });

      if (!response.ok) {
        throw new Error("Signup failed");
      }

      message.success("Signup successful!");
      navigate("/login");
    } catch (err) {
      console.error("Signup failed:", err);
      setError("Failed to sign up. Please try again.");
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
          Signup Page
        </Title>
        <Form layout="vertical" form={form} onFinish={submitSignupData}>
          <Form.Item
            label="Name"
            name="name"
            rules={[{ required: true, message: "Please enter your name!" }]}
          >
            <Input
              style={{ borderRadius: "4px" }}
              placeholder="Enter your name"
            />
          </Form.Item>

          <Form.Item
            label="Email"
            name="email"
            rules={[
              { required: true, message: "Please enter your email!" },
              { type: "email", message: "Please enter a valid email!" },
            ]}
          >
            <Input
              style={{ borderRadius: "4px" }}
              placeholder="Enter your email"
            />
          </Form.Item>

          <Form.Item
            label="Password"
            name="password"
            rules={[
              { required: true, message: "Please enter your password!" },
              {
                min: 6,
                message: "Password must be at least 6 characters long!",
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
              Signup
            </Button>
          </Form.Item>
        </Form>

        <Text style={{ textAlign: "center" }}>
          Already a member?{" "}
          <Link to="/login" style={{ color: "#1890ff" }}>
            Login
          </Link>{" "}
          here
        </Text>
      </Content>
    </Layout>
  );
}

export default SignupInputs;
