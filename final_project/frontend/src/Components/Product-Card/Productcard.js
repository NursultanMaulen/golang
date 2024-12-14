import React, { useState } from "react";
import { Card, Typography, Button, InputNumber, message } from "antd";
import { useLoginSignupContext } from "../../Context/LoginSignupContext";

const { Title, Paragraph, Text } = Typography;

function ProductCard({
  product = {},
  onAddToCart,
  onRemoveFromCart,
  onUpdateQuantity,
  showAddToCart = true,
  showRemoveButton = false,
  showQuantityInput = false,
  showInputNumber = true,
}) {
  const { name, description, id, price, quantity } = product;
  const { state: authState } = useLoginSignupContext();
  const { user } = authState;
  const [localQuantity, setLocalQuantity] = useState(quantity || 1);
  const [newQuantity, setNewQuantity] = useState(quantity);

  const handleAddToCart = async () => {
    try {
      if (!user?.user_id) {
        message.error("You need to log in to add products to your cart.");
        return;
      }
      console.log(user.user_id);
      const response = await fetch("http://localhost:8080/api/cart/add", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          user_id: user.user_id,
          product_id: id,
          quantity: localQuantity,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to add product to cart");
      }

      const result = await response.json();
      message.success(result.message || "Product added to cart successfully!");
    } catch (error) {
      console.error("Error adding product to cart:", error);
      message.error("Failed to add product to cart. Please try again.");
    }
  };

  const handleRemoveFromCart = async () => {
    if (onRemoveFromCart) {
      await onRemoveFromCart(id);
    }
  };

  const handleUpdateQuantity = async () => {
    if (!user?.user_id || user?.role !== "admin") {
      message.error(
        "You need to log in as an admin to edit product quantities."
      );
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/api/products/${id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Role: user.role,
        },
        body: JSON.stringify({ quantity: newQuantity }),
      });

      if (!response.ok) {
        throw new Error("Failed to update product quantity");
      }

      const result = await response.json();
      message.success(
        result.message || "Product quantity updated successfully!"
      );

      // Optional callback to update the UI
      onUpdateQuantity && onUpdateQuantity(id, newQuantity);
    } catch (error) {
      console.error("Error updating product quantity:", error);
      message.error("Failed to update product quantity. Please try again.");
    }
  };

  return (
    <Card
      hoverable
      style={{
        width: "100%",
        borderRadius: "8px",
        marginBottom: "16px",
        boxShadow: "0 4px 8px rgba(0, 0, 0, 0.2)",
      }}
    >
      <Title level={4}>{name}</Title>
      <Paragraph>{description}</Paragraph>
      <Text strong>{`Price: $${price.toFixed(2)} per unit`}</Text>

      {/* Show editable quantity if the user is an admin */}
      {user?.role === "admin" && (
        <div style={{ marginTop: "16px" }}>
          <InputNumber
            min={1}
            value={newQuantity}
            onChange={(value) => setNewQuantity(value)}
            style={{ maxWidth: "80px", marginRight: "8px" }}
          />
          <Button type="primary" onClick={handleUpdateQuantity}>
            Update Quantity
          </Button>
        </div>
      )}

      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          marginTop: "16px",
        }}
      >
        {showInputNumber && (
          <InputNumber
            min={1}
            value={localQuantity}
            onChange={(value) => setLocalQuantity(value)}
            style={{ maxWidth: "80px" }}
          />
        )}
        {showAddToCart && (
          <Button type="primary" onClick={handleAddToCart}>
            Add to Cart
          </Button>
        )}
        {showQuantityInput && <Text>{`Quantity: ${quantity}`}</Text>}
        {showRemoveButton && (
          <Button type="primary" onClick={handleRemoveFromCart}>
            Remove
          </Button>
        )}
      </div>
    </Card>
  );
}

export default ProductCard;
