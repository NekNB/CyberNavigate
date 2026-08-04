import { GetFile } from "../app/api/Simulator/Simulator";

export default class FileDownloader {
  // 1. private: скрыл массив внутри класса, чтобы его нельзя было изменить напрямую снаружи
  // 2. readonly: массив нельзя будет перезаписать (например, downloadedFiles = []),
  // но push() работать будет, так как readonly запрещает изменение ссылок, а не содержимого.
  private readonly downloadedFiles: string[] = [];

  /**
   * Имитация скачивания файла
   */
  public async downloadFile(fileId: string): Promise<boolean> {
    const file = await GetFile(fileId);
    this.downloadedFiles.push(fileId);
    return file.isSafe;
  }

  /**
   * Проверяет, был ли файл уже скачан
   */
  public isDownloaded(fileId: string): boolean {
    // Используем includes для проверки наличия элемента
    return this.downloadedFiles.includes(fileId);
  }

  /**
   * (Дополнительно) Геттер, если все же нужно получить список всех скачанных файлов снаружи
   * Возвращает копию массива, чтобы никто не смог сделать mutation (изменить) оригинальный массив
   */
  public getDownloadedList(): string[] {
    return [...this.downloadedFiles];
  }
}
