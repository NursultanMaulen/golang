import { NavLink as Link } from "react-router-dom";
import { Layout, Menu } from "antd";
import {
  HomeOutlined,
  SearchOutlined,
  ShoppingCartOutlined,
} from "@ant-design/icons";

const { Sider } = Layout;

function Sidebar() {
  return (
    <Sider
      width={150}
      style={{
        position: "relative",
        left: 0,
        background: "#fff",
        boxShadow: "2px 0 8px rgba(0, 0, 0, 0.1)",
        zIndex: 1000,
        overflow: "auto",
      }}
    >
      <Menu mode="inline" style={{ height: "100%", paddingBottom: "100px" }}>
        <Menu.Item key="1" icon={<HomeOutlined />}>
          <Link to="/">Home</Link>
        </Menu.Item>

        <Menu.Item key="2" icon={<SearchOutlined />}>
          <Link to="/explore">Explore</Link>
        </Menu.Item>

        <Menu.Item key="3" icon={<ShoppingCartOutlined />}>
          <Link to="/likes">Cart</Link>
        </Menu.Item>
      </Menu>
    </Sider>
  );
}

export default Sidebar;
