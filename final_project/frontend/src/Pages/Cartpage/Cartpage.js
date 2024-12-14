import React, { useEffect, useState } from "react";
import { Layout, Typography, Row, Col, Spin, Button, message } from "antd";
import { Footer, Header, Sidebar } from "../../Components/IndexAllComponents";
import ProductCard from "../../Components/Product-Card/Productcard";
import { useLoginSignupContext } from "../../Context/LoginSignupContext";

const { Content } = Layout;
const { Title, Text } = Typography;

function CartPage() {
  const { user } = useLoginSignupContext();
  const [cartItems, setCartItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [totalAmount, setTotalAmount] = useState(0);

  // Fetch Cart Items
  useEffect(() => {
    const fetchCartItems = async () => {
      if (!user?.user_id) {
        setCartItems([]);
        setLoading(false);
        return;
      }
      try {
        const response = await fetch(
          `http://localhost:8080/api/cart/items?user_id=${user.user_id}`
        );
        if (!response.ok) {
          throw new Error("Failed to fetch cart items");
        }
        const data = await response.json();
        setCartItems(data || []); // Default to an empty array if `data` is null or undefined

        // Calculate total amount
        const total =
          data && data.length > 0
            ? data.reduce((sum, item) => sum + item.price * item.quantity, 0)
            : 0;
        setTotalAmount(total);
      } catch (error) {
        console.error("Error fetching cart items:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchCartItems();
  }, [user?.user_id]);

  const handleOrder = async () => {
    if (cartItems.length === 0) {
      message.warning("Your cart is empty.");
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/api/orders", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          user_id: user.user_id,
          items: cartItems.map((item) => ({
            product_id: item.product_id,
            quantity: item.quantity,
            price: item.price,
          })),
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to place order");
      }

      const result = await response.json();
      message.success(result.message || "Order placed successfully!");
      setCartItems([]); 
      setTotalAmount(0); 
    } catch (error) {
      console.error("Error placing order:", error);
      message.error("Failed to place order. Please try again.");
    }
  };

  const handleRemoveFromCart = async (productId) => {
    try {
      const response = await fetch(
        `http://localhost:8080/api/cart/${productId}?user_id=${user.user_id}`,
        {
          method: "DELETE",
        }
      );

      if (!response.ok) {
        throw new Error("Failed to remove product from cart");
      }

      setCartItems((prevItems) =>
        prevItems.filter((item) => item.product_id !== productId)
      );

      // Recalculate the total amount
      const updatedTotal = cartItems
        .filter((item) => item.product_id !== productId)
        .reduce((sum, item) => sum + item.price * item.quantity, 0);
      setTotalAmount(updatedTotal);

      message.success("Product removed from cart successfully!");
    } catch (error) {
      console.error("Error removing product from cart:", error);
      message.error("Failed to remove product from cart. Please try again.");
    }
  };

  if (loading) {
    return (
      <div style={{ textAlign: "center", padding: "20px" }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <Layout style={{ minHeight: "100vh", overflow: "hidden" }}>
      <Header />
      <Layout>
        <Sidebar />
        <Content style={{ padding: "30px", minHeight: "100vh" }}>
          <Title
            level={2}
            style={{ textAlign: "center", marginBottom: "20px" }}
          >
            {cartItems.length > 0 ? "Your Cart" : "Your Cart is Empty"}
          </Title>
          <Row gutter={[16, 16]} justify="center">
            {cartItems.map((item) => (
              <Col key={item.product_id} xs={24} sm={12} md={8} lg={10}>
                <ProductCard
                  product={{ ...item }}
                  showRemoveButton={true}
                  showAddToCart={false}
                  onRemoveFromCart={() => handleRemoveFromCart(item.product_id)}
                  showQuantityInput={true}
                  showInputNumber={false}
                />
              </Col>
            ))}
          </Row>
          {cartItems.length > 0 && (
            <div
              style={{
                marginTop: "20px",
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
              }}
            >
              <Text strong style={{ fontSize: "18px" }}>
                Total Amount: ${totalAmount.toFixed(2)}
              </Text>
              <Button type="primary" onClick={handleOrder}>
                Place Order
              </Button>
            </div>
          )}
        </Content>
      </Layout>
      <Footer />
    </Layout>
  );
}

export default CartPage;
