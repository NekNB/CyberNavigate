import { useState, type FC } from "react";
import { RegisterUser } from "../../app/api/User/User";
import styles from "./Register.module.css";

interface RegisterProps {
  onClose: () => void;
  openLogin: () => void;
}

const Register: FC<RegisterProps> = ({ onClose, openLogin }) => {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [repeatedPassword, setRepeatedPassword] = useState("");
  const [errorMessage, setErrorMessage] = useState("");

  const isFormFilled =
    login !== "" &&
    password !== "" &&
    repeatedPassword !== "" &&
    repeatedPassword === password;

  // Проверяем, есть ли какая-то ошибка (для отображения в UI)
  const hasError = errorMessage !== "";

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div
        className={styles.registerWindow}
        onClick={(e) => e.stopPropagation()}
      >
        <h3 className={styles.title}>Регистрация</h3>

        <form
          className={styles.form}
          onSubmit={async (e) => {
            e.preventDefault();
            if (!isFormFilled) return;
            setErrorMessage(""); // Сбрасываем ошибку перед запросом
            const responseStatus = await RegisterUser(login, password);
            if (responseStatus === 400) {
              setErrorMessage("Пользователь с таким именем уже существует");
              return;
            }
            openLogin();
          }}
        >
          <input
            id="FrmLogin"
            type="text"
            placeholder="Логин"
            className={styles.input}
            onChange={(e) => setLogin(e.target.value.trim())}
          />
          <input
            type="password"
            placeholder="Пароль"
            className={styles.input}
            id="FrmPass"
            value={password} // Добавил value для полноты контроля
            onChange={(e) => setPassword(e.target.value.trim())}
          />
          <input
            type="password"
            id="FrmRepeatedPass"
            placeholder="Повторите Пароль"
            className={styles.input}
            value={repeatedPassword} // Добавил value
            onChange={(e) => {
              const val = e.target.value.trim();
              setRepeatedPassword(val);
              if (password !== val && val !== "") {
                setErrorMessage("Пароли не совпадают");
              } else {
                setErrorMessage(""); // Сбрасываем, если исправили
              }
            }}
          />

          {hasError ? <p className={styles.errors}>{errorMessage}</p> : null}

          <div className={styles.buttonGroup}>
            <button
              type="submit"
              className={`${styles.button} ${
                isFormFilled ? styles.buttonActive : styles.buttonDisabled
              }`}
              disabled={!isFormFilled}
            >
              Зарегистрироваться
            </button>

            <button
              type="button"
              className={styles.buttonSecondary}
              onClick={openLogin}
            >
              Назад
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default Register;
