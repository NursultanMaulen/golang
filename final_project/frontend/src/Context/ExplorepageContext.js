import {
  createContext,
  useContext,
  useEffect,
  useReducer,
  useCallback,
  useMemo,
} from "react";

const ExplorePageContext = createContext();
export const useExplorePageContext = () => useContext(ExplorePageContext);

function ExplorepageContext({ children }) {
  const initialState = {
    products: [],
    isLoading: false,
  };

  function reducerFn(state, action) {
    switch (action.type) {
      case "FETCH_PRODUCTS_SUCCESS":
        return { ...state, products: action.payload, isLoading: false };
      case "FETCH_PRODUCTS_LOADING":
        return { ...state, isLoading: true };
      case "FETCH_PRODUCTS_ERROR":
        return { ...state, isLoading: false };
      default:
        return state;
    }
  }

  const [state, dispatch] = useReducer(reducerFn, initialState);

  const fetchProducts = useCallback(async () => {
    dispatch({ type: "FETCH_PRODUCTS_LOADING" });
    try {
      const response = await fetch("http://localhost:8080/api/products");
      if (!response.ok) {
        throw new Error("Failed to fetch products");
      }
      const data = await response.json();

      const normalizedData = data.map((product) => ({
        id: product.product_id || product.id,
        name: product.name || "Unknown Product",
        description: product.description || "No description available",
        price: product.price || 0, 
      }));

      dispatch({ type: "FETCH_PRODUCTS_SUCCESS", payload: normalizedData });
    } catch (error) {
      console.error("Error fetching products:", error);
      dispatch({ type: "FETCH_PRODUCTS_ERROR" });
    }
  }, []);

  useEffect(() => {
    fetchProducts();
  }, [fetchProducts]);

  return (
    <ExplorePageContext.Provider value={{ state, fetchProducts }}>
      {children}
    </ExplorePageContext.Provider>
  );
}

export default ExplorepageContext;
