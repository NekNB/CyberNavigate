import axios from "axios";
import { loadConfig } from "../config";

export const apiServer = axios.create({
  timeout: 10000,
  baseURL: loadConfig().internalUrl,
  headers: {
    "Content-Type": "application/json",
  },
});
