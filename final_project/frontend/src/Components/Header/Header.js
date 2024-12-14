import { Link, useNavigate } from "react-router-dom";
import { useExplorePageContext } from "../../Context/IndexAllContext";
import { Layout, AutoComplete, Input, Button } from "antd";
import { SearchOutlined } from "@ant-design/icons";
import { MdLogin, MdAccountCircle } from "react-icons/md";
import { useState } from "react";

const { Header: AntHeader } = Layout;

function Header() {
  const { dispatch } = useExplorePageContext();
  const [searchTerm, setSearchTerm] = useState("");
  const token = localStorage.getItem("token");
  const navigate = useNavigate();

  const options = [{ value: "home" }, { value: "explore" }, { value: "cart" }];

  const handleSearch = () => {
    if (searchTerm.trim().toLowerCase() === "home") {
      navigate("/");
    } else if (searchTerm.trim().toLowerCase() === "explore") {
      navigate("/explore");
    } else if (searchTerm.trim().toLowerCase() === "cart") {
      navigate("/cart");
    } else if (searchTerm.trim().toLowerCase() === "account") {
      navigate("/accounts");
    } else if (searchTerm.trim()) {
      navigate(`/search?query=${searchTerm}`);
    }
  };

  return (
    <AntHeader
      style={{
        background: "#fff",
        padding: 0,
        borderBottom: "1px solid #d9d9d9",
        boxShadow: "0 1px 4px rgba(0, 0, 0, 0.1)",
      }}
    >
      <div
        style={{
          display: "flex",
          alignItems: "center",
          height: "64px",
          padding: "0 20px",
        }}
      >
        <Link to="/" style={{ flexShrink: 0 }}>
          <div style={{ display: "flex", alignItems: "center" }}>
            <img
              src={require("../../assets/delivery-box.png")}
              alt="logo"
              style={{ marginRight: "8px", height: "40px", width: "40px" }}
            />
            <span
              style={{ fontSize: "24px", fontWeight: "bold", color: "purple" }}
            >
              Stream
            </span>
            <span
              style={{ fontSize: "24px", fontWeight: "bold", color: "red" }}
            >
              Box
            </span>
          </div>
        </Link>

        <div style={{ flex: 1, margin: "0 20px", display: "flex" }}>
          <AutoComplete
            options={options}
            style={{ maxWidth: "600px", width: "100%" }}
            value={searchTerm}
            onChange={(value) => setSearchTerm(value)}
            onSelect={(value) => setSearchTerm(value)}
          >
            <Input
              placeholder="Search item"
              suffix={<SearchOutlined onClick={handleSearch} />}
              onPressEnter={handleSearch}
            />
          </AutoComplete>
        </div>

        <div style={{ flexShrink: 0 }}>
          {!token ? (
            <Link to="/login">
              <Button type="primary" icon={<MdLogin />}>
                Login
              </Button>
            </Link>
          ) : (
            <Link to="/accounts">
              <Button type="primary" icon={<MdAccountCircle />}>
                Account
              </Button>
            </Link>
          )}
        </div>
      </div>
    </AntHeader>
  );
}

export default Header;
