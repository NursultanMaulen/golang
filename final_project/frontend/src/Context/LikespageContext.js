import { createContext, useContext, useReducer } from "react";

const likeContext = createContext();
export const useLikeContext = () => useContext(likeContext);

function LikespageContext({ children }) {
  function reducerFn(state, action) {
    switch (action.type) {
      case "ADD_TO_LIKES":
        return {
          ...state,
          likedProducts: [...state.likedProducts, action.payload],
        };
      case "REMOVE_FROM_LIKES":
        return {
          ...state,
          likedProducts: state.likedProducts.filter(
            (product) => product.id !== action.payload.id
          ),
        };
      default:
        return state;
    }
  }

  const [state, dispatch] = useReducer(reducerFn, {
    likedProducts: [], // Список лайкнутых продуктов
  });

  const { likedProducts } = state;

  return (
    <likeContext.Provider
      value={{
        likedProducts,
        dispatch,
      }}
    >
      {children}
    </likeContext.Provider>
  );
}

export default LikespageContext;
