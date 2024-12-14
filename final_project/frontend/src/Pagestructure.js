import App from "./App";
import {
  ExplorepageContext,
  LikespageContext,
  LoginSignupContext,
} from "./Context/CoreContextFiles";
function Pagestructure() {
  return (
    <div>
      <ExplorepageContext>
        <LikespageContext>
          <LoginSignupContext>
            <App />
          </LoginSignupContext>
        </LikespageContext>
      </ExplorepageContext>
    </div>
  );
}

export default Pagestructure;
