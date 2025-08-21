function main() {
  const lines = parse();

  // pbl, ptl, pbr, ptr := Point{7, 7, 0}, Point{7, 27, 0}, Point{27, 7, 0}, Point{27, 27, 0}

  // top, bottom, left, right := Line{ptl, ptr}, Line{pbl, pbr}, Line{pbl, ptl}, Line{pbr, ptr}

  const [minX, minY, maxX, maxY] = [7, 7, 27, 27];

  const intersections: Point[] = [];

  for (
    let [i, ...nextLines] = lines;
    nextLines.length > 0;
    [i, ...nextLines] = nextLines
  ) {
    const line = i.line;

    for (const next of nextLines) {
      const l = next.line;

      const [intersectionPoint, intersects] = lineIntersection(line, l);
      if (intersects && intersectionPoint) {
        // fmt.Println("Hailstone A:", fmt.Sprintf("%f, %f, %f @ %f, %f, %f", line.p1.x, line.p1.y, line.p1.z, line.p2.x, line.p2.y, line.p2.z))
        // fmt.Println("Hailstone B:", fmt.Sprintf("%f, %f, %f @ %f, %f, %f", l.p1.x, l.p1.y, l.p1.z, l.p2.x, l.p2.y, l.p2.z))
        // fmt.Println("Hailstone ':", fmt.Sprintf("paths will cross at (x=%f,y=%f)", intersectionPoint.x, intersectionPoint.y))
        // fmt.Println()

        const nanoSecond = getNano(i, intersectionPoint);
        // fmt.Println(nanoSecond)
        if (
          intersectionPoint.x >= minX &&
          intersectionPoint.x <= maxX &&
          intersectionPoint.y >= minY &&
          intersectionPoint.y <= maxY &&
          nanoSecond >= 0
        ) {
          // fmt.Println()
          intersections.push(intersectionPoint);
        }
      }
    }
  }

  console.log(intersections);
}

// 19 - 2 - 2 = 15
// x = 19 nx = 16 vx = -2  / s
// x + (s * vx) = nx
// s * vx = nx - x
// s = (nx - x) / vx

function getNano(l: LineVelocity, p: Point): number {
  return (p.x - l.line.p1.x) / l.velocity.x;
}

function parse(): LineVelocity[] {
  const lines: LineVelocity[] = [];

  const text = await Deno.readTextFile('assets/2023/task23.txt');
  text.split('\n').foreach((line: string) => {
    const [point, velocity] = line.split(' @ ');
    const pointCoords = point.split(', ');
    const velocityValues = velocity.split(', ');

    const [x, y, z] = [
      parseInt(pointCoords[0]),
      parseInt(pointCoords[1]),
      parseInt(pointCoords[2]),
    ];
    const [vx, vy, vz] = [
      parseInt(trimWhitespace(velocityValues[0])),
      parseInt(trimWhitespace(velocityValues[1])),
      parseInt(trimWhitespace(velocityValues[2])),
    ];

    const p1: Point = { x, y, z };
    const p2: Point = { x: x + vx, y: y + vy, z: z + vz };

    lines.push({ line: { p1, p2 }, velocity: { x: vx, y: vy, z: vz } });
  });

  return lines;
}

function trimWhitespace(s: string): string {
  return s.trim();
}

function lineIntersection(l1: Line, l2: Line): [Point | undefined, boolean] {
  const a1 = l1.p1.y - l1.p2.y;
  const b1 = l1.p1.x - l1.p2.x;
  const c1 = l1.p1.x * l1.p2.y - l1.p2.x * l1.p1.y;

  const a2 = l2.p1.y - l2.p2.y;
  const b2 = l2.p1.x - l2.p2.x;
  const c2 = l2.p1.x * l2.p2.y - l2.p2.x * l2.p1.y;

  const px_n = c1 * b2 - b1 * c2;
  const py_n = c1 * a2 - a1 * c2;
  const denominator = b1 * a2 - a1 * b2;

  if (denominator == 0) {
    return [undefined, false];
  }

  return [{ x: px_n / denominator, y: py_n / denominator, z: 0 }, true];
}

function segmentIntersection(
  l1: Segment,
  l2: Segment,
): [Point | undefined, boolean] {
  const [P1, P2, P3, P4] = [l1.start, l1.end, l2.start, l2.end];
  const [Ax, Ay] = [P2.x - P1.x, P2.y - P1.y];
  const [Bx, By] = [P3.x - P4.x, P3.y - P4.y];
  const [Cx, Cy] = [P1.x - P3.x, P1.y - P3.y];

  const alphaNumerator = By * Cx - Bx * Cy;
  const betaNumerator = Ax * Cy - Ay * Cx;
  const denominator = Ay * Bx - Ax * By;

  if (denominator == 0) {
    return [
      undefined,
      alphaNumerator == 0 &&
        betaNumerator == 0 &&
        onLineX(P1, P2, P3, P4) &&
        onLineY(P1, P2, P3, P4),
    ];
  }

  const [px, py] = [alphaNumerator / denominator, betaNumerator / denominator];

  return [{ x: px, y: py, z: 0 }, between(px, 0, 1) && between(py, 0, 1)];
}

function between(n: number, start: number, end: number): boolean {
  return n >= start && n <= end;
}

function onLineX(P1: Point, P2: Point, P3: Point, P4: Point): boolean {
  return (
    (P1.x >= P3.x && P1.x <= P4.x) ||
    (P2.x >= P3.x && P2.x <= P4.x) ||
    (P3.x >= P1.x && P3.x <= P2.x) ||
    (P4.x >= P1.x && P4.x <= P2.x)
  );
}

function onLineY(P1: Point, P2: Point, P3: Point, P4: Point): boolean {
  return (
    (P1.y >= P3.y && P1.y <= P4.y) ||
    (P2.y >= P3.y && P2.y <= P4.y) ||
    (P3.y >= P1.y && P3.y <= P2.y) ||
    (P4.y >= P1.y && P4.y <= P2.y)
  );
}

type Point = {
  x: number;
  y: number;
  z: number;
};

type Line = {
  p1: Point;
  p2: Point;
};

type LineVelocity = {
  line: Line;
  velocity: Point;
};

type Segment = {
  start: Point;
  end: Point;
};
