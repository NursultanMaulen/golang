import { useExplorePageContext } from "../../Context/ExplorepageContext";
import { Button } from "antd";
import { useLoginSignupContext } from "../../Context/LoginSignupContext";

function AddProductButton() {
  const { createProduct } = useExplorePageContext();
  const { state: authState } = useLoginSignupContext();
  const { isAuthenticated } = authState;

  const handleAddProduct = () => {
    const newProduct = {
      title: "ADDED Product",
      description: "Description of the new Product",
      category: "Frontend Development",
      thumbnailUrl: "https://example.com/thumbnail.jpg",
      ProductUrl: "https://www.youtube.com/embed/TBIjgBVFjVI",
      creator_pic:
        "https://yt3.googleusercontent.com/ytc/AIdro_mKzklyPPhghBJQH5H3HpZ108YcE618DBRLAvRUD1AjKNw=s160-c-k-c0x00ffffff-no-rj",
      creator_name: "FireShip",
    };

    createProduct(newProduct);
  };

  return (
    <Button disabled={!isAuthenticated} onClick={handleAddProduct}>
      Add Product
    </Button>
  );
}

export default AddProductButton;
