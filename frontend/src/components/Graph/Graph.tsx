// Graph.tsx
import React, { useCallback, useEffect, useMemo, useRef } from "react";
import styles from "./Graph.module.css";

interface GraphProps {
  /** Массив значений от -100 до 100 */
  data: number[];
  /** Отступы от краёв (по умолчанию 40) */
  padding?: number;
  /** Количество интерполированных точек между узлами (по умолчанию 25) */
  segments?: number;
  /** Цвет нулевой линии */
  zeroLineColor?: string;
  /** Дополнительный класс для контейнера */
  className?: string;
  /** Минимальное расстояние между точками по X (по умолчанию 15px) */
  minPointSpacing?: number;
  /** Соотношение сторон (ширина/высота), по умолчанию 3:1 */
  aspectRatio?: number;
}

interface Point {
  x: number;
  y: number;
}

const Graph: React.FC<GraphProps> = ({
  data,
  padding = 40,
  segments = 25,
  zeroLineColor = "rgba(255, 255, 255, 0.7)",
  className = "",
  minPointSpacing = 15,
  aspectRatio = 3,
}) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  // Вычисляем размеры на основе данных и контейнера
  const dimensions = useMemo(() => {
    const pointCount = data.length;
    // Минимальная ширина графика = отступы + расстояние между точками * (количество точек - 1)
    const minGraphWidth = padding * 2 + minPointSpacing * (pointCount - 1);

    // Высота на основе соотношения сторон
    const height = minGraphWidth / aspectRatio;

    return {
      width: minGraphWidth,
      height: Math.max(height, 200), // Минимальная высота 200px
    };
  }, [data.length, padding, minPointSpacing, aspectRatio]);

  const width = dimensions.width;
  const height = dimensions.height;
  const graphWidth = width - padding * 2;
  const graphHeight = height - padding * 2;
  const zeroY = padding + graphHeight / 2;
  const scaleY = graphHeight / 200; // диапазон -100..100 = 200 единиц

  /**
   * Преобразование значения данных в координаты на canvas
   */
  const valueToCoords = useCallback(
    (index: number, value: number): Point => {
      const x = padding + (index / (data.length - 1)) * graphWidth;
      const y = zeroY - value * scaleY;
      return { x, y };
    },
    [data.length, padding, graphWidth, zeroY, scaleY],
  );

  /**
   * Кубическая интерполяция Catmull-Rom
   */
  const catmullRom = useCallback(
    (p0: Point, p1: Point, p2: Point, p3: Point, t: number): Point => {
      const t2 = t * t;
      const t3 = t2 * t;

      const x =
        0.5 *
        (2 * p1.x +
          (-p0.x + p2.x) * t +
          (2 * p0.x - 5 * p1.x + 4 * p2.x - p3.x) * t2 +
          (-p0.x + 3 * p1.x - 3 * p2.x + p3.x) * t3);

      const y =
        0.5 *
        (2 * p1.y +
          (-p0.y + p2.y) * t +
          (2 * p0.y - 5 * p1.y + 4 * p2.y - p3.y) * t2 +
          (-p0.y + 3 * p1.y - 3 * p2.y + p3.y) * t3);

      return { x, y };
    },
    [],
  );

  /**
   * Получение массива точек плавной кривой
   */
  const getCurvePoints = useCallback((): Point[] => {
    const points: Point[] = [];
    const n = data.length;

    if (n < 2) return points;

    const nodePoints = data.map((val, i) => valueToCoords(i, val));

    for (let i = 0; i < n - 1; i++) {
      const p0 = nodePoints[Math.max(0, i - 1)];
      const p1 = nodePoints[i];
      const p2 = nodePoints[i + 1];
      const p3 = nodePoints[Math.min(n - 1, i + 2)];

      for (let j = 0; j <= segments; j++) {
        const t = j / segments;
        points.push(catmullRom(p0, p1, p2, p3, t));
      }
    }

    return points;
  }, [data, valueToCoords, catmullRom, segments]);

  /**
   * Отрисовка нулевой линии
   */
  const drawZeroLine = useCallback(
    (ctx: CanvasRenderingContext2D) => {
      ctx.strokeStyle = zeroLineColor;
      ctx.lineWidth = 2;
      ctx.setLineDash([]);

      ctx.beginPath();
      ctx.moveTo(padding, zeroY);
      ctx.lineTo(width - padding, zeroY);
      ctx.stroke();

      ctx.fillStyle = "rgba(255, 255, 255, 0.9)";
      ctx.font = "bold 12px Arial";
      ctx.textAlign = "right";
      ctx.fillText("", padding - 10, zeroY + 4);
    },
    [zeroLineColor, padding, width, zeroY],
  );

  /**
   * Отрисовка плавной линии данных
   */
  const drawDataLine = useCallback(
    (ctx: CanvasRenderingContext2D) => {
      const points = getCurvePoints();
      if (points.length < 2) return;

      // Свечение под линией
      ctx.strokeStyle = "rgba(255, 255, 255, 0.15)";
      ctx.lineWidth = 8;
      ctx.lineJoin = "round";
      ctx.lineCap = "round";
      ctx.setLineDash([]);
      ctx.beginPath();
      ctx.moveTo(points[0].x, points[0].y);
      for (let i = 1; i < points.length; i++) {
        ctx.lineTo(points[i].x, points[i].y);
      }
      ctx.stroke();

      // Основная линия
      ctx.strokeStyle = "#2CC295";
      ctx.lineWidth = 5;
      ctx.beginPath();
      ctx.moveTo(points[0].x, points[0].y);
      for (let i = 1; i < points.length; i++) {
        ctx.lineTo(points[i].x, points[i].y);
      }
      ctx.stroke();
    },
    [getCurvePoints],
  );

  /**
   * Основная функция отрисовки
   */
  const draw = useCallback(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    // Очистка canvas
    ctx.clearRect(0, 0, width, height);

    // Отрисовка всех слоёв
    drawZeroLine(ctx);
    drawDataLine(ctx);
  }, [width, height, drawZeroLine, drawDataLine]);

  // Перерисовка при изменении данных или размеров
  useEffect(() => {
    draw();
  }, [draw]);

  // Настройка canvas и обработка ресайза
  useEffect(() => {
    const canvas = canvasRef.current;
    const container = containerRef.current;
    if (!canvas || !container) return;

    const updateCanvasSize = () => {
      const dpr = window.devicePixelRatio || 1;

      // Получаем доступную ширину контейнера
      const containerWidth = container.clientWidth;

      // Если контейнер уже, чем минимальная ширина графика, используем минимальную
      // Если шире — растягиваем график на всю ширину контейнера
      const displayWidth = Math.max(containerWidth, dimensions.width);
      const scaleFactor = displayWidth / dimensions.width;
      const displayHeight = dimensions.height * scaleFactor;

      // Устанавливаем размеры canvas с учётом DPR
      canvas.width = displayWidth * dpr;
      canvas.height = displayHeight * dpr;
      canvas.style.width = `${displayWidth}px`;
      canvas.style.height = `${displayHeight}px`;

      const ctx = canvas.getContext("2d");
      if (ctx) {
        ctx.setTransform(dpr * scaleFactor, 0, 0, dpr * scaleFactor, 0, 0);
      }
    };

    // Создаём ResizeObserver для отслеживания изменений размера контейнера
    const resizeObserver = new ResizeObserver(() => {
      updateCanvasSize();
      draw();
    });

    resizeObserver.observe(container);
    updateCanvasSize();

    return () => {
      resizeObserver.disconnect();
    };
  }, [dimensions, draw]);

  const canvasClassName = [styles.canvas].filter(Boolean).join(" ");

  return (
    <div
      ref={containerRef}
      className={`${styles.container} ${className}`.trim()}
      style={{ width: "100%" }}
    >
      <canvas ref={canvasRef} className={canvasClassName} />
    </div>
  );
};

export default Graph;
