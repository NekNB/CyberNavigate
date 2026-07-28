import { useState, type FC } from "react";
import { Login as LoginUser } from "../../app/api/Auth/Auth";
import styles from "./Login.module.css";

interface LoginProps {
  onClose: () => void;
  openRegister: () => void;
}

const Login: FC<LoginProps> = ({ onClose, openRegister }) => {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [errorMessage, setErrorMessage] = useState("");

  // trim() убирает пробелы в начале и конце
  const isFormFilled = login.trim() !== "" && password.trim() !== "";

  // Проверяем, есть ли какая-то ошибка (для отображения в UI)
  const hasError = errorMessage !== "";
  return (
    <div className={styles.overlay} onClick={onClose}>
      <div
        className={styles.registerWindow}
        onClick={(e) => e.stopPropagation()}
      >
        <h3 className={styles.title}>Вход</h3>

        <form
          className={styles.form}
          onSubmit={async (e) => {
            e.preventDefault();
            if (!isFormFilled) return; // Дополнительная защита
            setErrorMessage(""); // Сбрасываем ошибку перед запросом
            const responseStatus = await LoginUser(login, password);
            if (responseStatus === 401) {
              setErrorMessage("Неверные логин или пароль");
              console.log("Неверные логин");
              return;
            }
            onClose();
          }}
        >
          {/* Привязываем значение поля к состоянию */}
          <input
            type="text"
            placeholder="Логин"
            className={styles.input}
            onChange={(e) => setLogin(e.target.value)}
          />
          <input
            type="password"
            placeholder="Пароль"
            className={styles.input}
            onChange={(e) => setPassword(e.target.value)}
          />

          {hasError ? <p className={styles.errors}>{errorMessage}</p> : null}
          <div className={styles.buttonGroup}>
            {/* Динамически меняем класс в зависимости от isFormFilled */}
            <button
              type="submit"
              className={`${styles.button} ${
                isFormFilled ? styles.buttonActive : styles.buttonDisabled
              }`}
              disabled={!isFormFilled} // Блокируем кнопку, если форма пустая
            >
              Войти
            </button>

            <button
              type="button"
              className={styles.buttonSecondary}
              onClick={openRegister}
            >
              Зарегистрироваться
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default Login;
