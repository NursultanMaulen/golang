import React, { useEffect, useState, useMemo, useCallback } from "react";
import { Table, Button, Typography, Card, Form, Input, message } from "antd";
import { useNavigate } from "../../Utils/CustomUtils";
import { useLoginSignupContext } from "../../Context/IndexAllContext";
import { logoutHandler } from "../../services/LoginSignUpServices";

const { Title } = Typography;

function Accountdetails() {
  const { loginData, dispatch } = useLoginSignupContext();
  const navigate = useNavigate();
  const [userData, setUserData] = useState({
    id: "",
    username: "",
    email: "",
    password_hash: "",
  });
  const [form] = Form.useForm();

  const fetchUserData = async (userId) => {
    try {
      const response = await fetch(`http://localhost:8080/api/users/${userId}`);
      if (!response.ok) {
        throw new Error("Failed to fetch user data");
      }

      const data = await response.json();
      setUserData({
        id: data.user_id,
        username: data.username,
        email: data.email,
        password_hash: data.password_hash,
      });
      localStorage.setItem("userData", JSON.stringify(data));
    } catch (error) {
      console.error("Error fetching user data:", error);
    }
  };

  useEffect(() => {
    const storedUserData = JSON.parse(localStorage.getItem("userData"));

    if (loginData && loginData.user_id) {
      localStorage.setItem("userData", JSON.stringify(loginData));
      updateUserDataState(loginData);
      fetchUserData(loginData.user_id);
    } else if (storedUserData) {
      updateUserDataState(storedUserData);
      fetchUserData(storedUserData.user_id);
    }
  }, [loginData]);

  const updateUserDataState = useCallback(
    (newUserData) => {
      setUserData((prevUserData) => {
        if (
          prevUserData.username === newUserData.username &&
          prevUserData.email === newUserData.email &&
          prevUserData.password_hash === newUserData.password_hash
        ) {
          return prevUserData;
        }
        return newUserData;
      });
      form.setFieldsValue(newUserData);
    },
    [form]
  );

  const logOutUserFromApp = useCallback(() => {
    logoutHandler(dispatch);
    localStorage.removeItem("userData");
    navigate("/login");
  }, [navigate]);

  const dataSource = useMemo(() => {
    return [
      {
        key: "1",
        username: userData.username,
        email: userData.email,
      },
    ];
  }, [userData]);

  const columns = [
    {
      title: "Name",
      dataIndex: "username",
      key: "username",
    },
    {
      title: "Email",
      dataIndex: "email",
      key: "email",
    },
    {
      title: "Logout",
      key: "logout",
      render: () => (
        <Button type="primary" danger onClick={logOutUserFromApp}>
          Logout
        </Button>
      ),
    },
  ];

  const updateUserData = async (values) => {
    const updatedUserData = {
      ...userData,
      username: values.username,
      email: values.email,
      password_hash: values.password_hash,
    };

    try {
      const response = await fetch(
        `http://localhost:8080/api/users/${userData.id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(updatedUserData),
        }
      );

      if (!response.ok) {
        throw new Error("Failed to update user");
      }

      await fetchUserData(userData.id);
      message.success("User data updated successfully!");
    } catch (error) {
      console.error("Error updating user data:", error);
      message.error("Failed to update user data. Please try again.");
    }
  };

  const deleteAccount = async () => {
    try {
      const response = await fetch(
        `http://localhost:8080/api/users/${userData.id}`,
        {
          method: "DELETE",
        }
      );

      if (!response.ok) {
        throw new Error("Failed to delete user");
      }

      logOutUserFromApp();
      message.success("Account deleted successfully!");
    } catch (error) {
      console.error("Error deleting user:", error);
      message.error("Failed to delete account. Please try again.");
    }
  };

  return (
    <Card style={{ maxWidth: "1000px", margin: "auto", padding: "24px" }}>
      <Title level={1} style={{ textAlign: "center" }}>
        Account Details
      </Title>
      <Table dataSource={dataSource} columns={columns} pagination={false} />

      <Form
        form={form}
        layout="vertical"
        onFinish={updateUserData}
        initialValues={userData}
        style={{ marginTop: "24px" }}
      >
        <Form.Item
          label="Name"
          name="username"
          rules={[{ required: true, message: "Please enter your name!" }]}
        >
          <Input placeholder="Enter your name" />
        </Form.Item>

        <Form.Item
          label="Email"
          name="email"
          rules={[{ required: true, message: "Please enter your email!" }]}
        >
          <Input type="email" placeholder="Enter your email" />
        </Form.Item>

        <Form.Item
          label="Password"
          name="password_hash"
          rules={[{ required: true, message: "Please enter your password!" }]}
        >
          <Input.Password placeholder="Enter your password" />
        </Form.Item>

        <Button type="primary" htmlType="submit">
          Update User Data
        </Button>
      </Form>

      <Button
        type="primary"
        danger
        style={{ marginTop: "16px" }}
        onClick={deleteAccount}
      >
        Delete Account
      </Button>
    </Card>
  );
}

export default Accountdetails;
