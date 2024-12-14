import React, { useEffect, useState } from "react";
import { useExplorePageContext } from "../../Context/ExplorepageContext";
import Productcard from "../../Components/Product-Card/Productcard";
import Sidebar from "../../Components/Sidebar/Sidebar";
import Header from "../../Components/Header/Header";
import Footer from "../../Components/Footer/Footer";
import { Layout, Spin, Row, Col, Pagination } from "antd";
import { useLoginSignupContext } from "../../Context/IndexAllContext";

const { Content } = Layout;

function ExplorePage() {
  const { state, fetchProducts } = useExplorePageContext();
  const { products, isLoading } = state;

  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(6);

  useEffect(() => {
    fetchProducts();
  }, [fetchProducts]);

  const handlePageChange = (page) => {
    setCurrentPage(page);
  };

  const paginatedProducts = products.slice(
    (currentPage - 1) * pageSize,
    currentPage * pageSize
  );

  if (isLoading) {
    return (
      <div style={{ textAlign: "center", padding: "20px" }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <Layout style={{ minHeight: "100vh", overflow: "hidden" }}>
      <Header />
      <Layout style={{ background: "#fff" }}>
        <Sidebar />
        <Content style={{ padding: "30px", minHeight: "100vh" }}>
          <Row gutter={[10, 10]} justify="center">
            {paginatedProducts.map((product) => (
              <Col key={product.id} xs={3} sm={2} md={8} lg={10}>
                <Productcard product={product} />
              </Col>
            ))}
          </Row>
          <Pagination
            current={currentPage}
            pageSize={pageSize}
            total={products.length}
            onChange={handlePageChange}
            showSizeChanger={false}
          />
        </Content>
        <Footer />
      </Layout>
    </Layout>
  );
}

export default ExplorePage;
