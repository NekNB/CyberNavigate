import { useEffect, useState, type FC } from "react";
import { useNavigate } from "react-router";
import { GetUser } from "../../../app/api/User/User";
import Login from "../../../components/Login/Login";
import ProgressBar from "../../../components/ProgressBar/ProgressBar";
import Register from "../../../components/Register/Register";
import styles from "./Loader.module.css";

export interface ILoaderProps {
  actions: { action: () => Promise<void>; placeholder?: string }[];
  setIsLoading: () => Promise<void>;
}

const Loader: FC<ILoaderProps> = ({ actions, setIsLoading }) => {
  const [progress, setProgress] = useState(0);
  const [isLoginOpen, setIsLoginOpen] = useState(false);
  const [isRegisterOpen, setIsRegisterOpen] = useState(false);
  const [task, setTask] = useState("Пока ждем");
  const navigate = useNavigate();

  useEffect(() => {
    const runActions = async () => {
      for (let i = 0; i < actions.length; i++) {
        setTask(actions[i].placeholder || "");
        await actions[i].action();
        setProgress(((i + 1) / actions.length) * 100);
      }
      await setIsLoading();
    };

    runActions();
  }, [actions]);

  useEffect(() => {
    const checkAuth = async () => {
      try {
        await GetUser();
        setIsLoginOpen(false);
      } catch (error) {
        setIsLoginOpen(true);
      }
    };
    checkAuth();
  }, [isLoginOpen]);

  return (
    <div className={styles.background}>
      {isLoginOpen && (
        <Login
          onClose={() => {
            navigate("/simulator");
            setIsLoginOpen(false);
          }}
          openRegister={() => {
            setIsLoginOpen(false);
            setIsRegisterOpen(true);
          }}
        />
      )}
      {isRegisterOpen && (
        <Register
          onClose={() => {
            setIsRegisterOpen(false);
          }}
          openLogin={() => {
            setIsRegisterOpen(false);
            setIsLoginOpen(true);
          }}
        />
      )}
      <div className={styles.loaderWrapper}>
        <img src="/assets/logo.svg" className={styles.logo} />
        <ProgressBar progress={progress} />
        <p>{task}</p>
      </div>
    </div>
  );
};

export default Loader;
