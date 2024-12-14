import { useExplorePageContext } from "../../Context/ExplorepageContext";
import { Button } from "antd";
import { DeleteOutlined } from "@ant-design/icons";
import { useLoginSignupContext } from "../../Context/LoginSignupContext";

function DeleteProductButton({ ProductId }) {
  const { deleteProduct } = useExplorePageContext();
  const { state: authState } = useLoginSignupContext();
  const { isAuthenticated } = authState;

  return (
    <Button
      icon={<DeleteOutlined />}
      danger
      onClick={() => deleteProduct(ProductId)}
      disabled={!isAuthenticated}
    ></Button>
  );
}

export default DeleteProductButton;
