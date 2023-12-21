package solutions.year2022;

import common.Problem;

import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;
import java.util.TreeMap;
import java.util.function.Function;
import java.util.stream.Collectors;
import java.util.stream.LongStream;
import java.util.stream.Stream;

@SuppressWarnings("unused")
public class Day09 extends Problem {

    public static final String file = "2022/09.txt";
    public static final Problem INSTANCE = new Day09(file, SolutionMethod.ALL);

    public static final double frequency = 0.5;
    public static final boolean visualize = false;

    public Day09(String file) {
        super(file);
    }

    public Day09(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    @Override
    public String solutionPartOne(Stream<String> data) {
        final var bridge = new Bridge(visualize);

        makeMoves(bridge, data);

        return String.valueOf(bridge.tail.tailMoves.size());
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        final var bridge = new Bridge(9, visualize);

        makeMoves(bridge, data);

        return String.valueOf(bridge.tail.tailMoves.size());
    }

    private void makeMoves(Bridge bridge, Stream<String> data) {
        data.map(move -> move.split(" "))
                .map(move -> new Move(Direction.of(move[0]), Long.parseLong(move[1])))
                .forEach(move -> bridge.move(move.direction(), move.steps()));
    }

    private enum Direction {
        UP("U"),
        DOWN("D"),
        LEFT("L"),
        RIGHT("R"),
        DOWN_RIGHT("DR"),
        DOWN_LEFT("DL"),
        UP_RIGHT("UR"),
        UP_LEFT("UL");

        final String value;

        Direction(String value) {
            this.value = value;
        }

        public static Direction of(Direction primary, Direction secondary) {
            if (primary == secondary) {
                return primary;
            }

            if (primary.isVertical() && secondary.isHorizontal()) {
                return Direction.of("%s%s".formatted(primary.getValue(), secondary.getValue()));
            }

            if (primary.isHorizontal() && secondary.isVertical()) {
                return Direction.of("%s%s".formatted(secondary.getValue(), primary.getValue()));
            }

            throw new IllegalArgumentException();
        }

        public static Direction of(String value) {
            return Arrays.stream(Direction.values())
                    .filter(direction -> direction.getValue().equals(value))
                    .findFirst()
                    .orElse(null);
        }

        public String getValue() {
            return value;
        }

        private Boolean isVertical() {
            return this == Direction.DOWN || this == Direction.UP;
        }

        private Boolean isHorizontal() {
            return this == Direction.LEFT || this == Direction.RIGHT;
        }
    }

    private record Rope(Position position, Map<Direction, Rope> next) {
    }

    private record Move(Direction direction, Long steps) {
    }

    private record Position(Long row, Long col) {

        public Position in(Direction direction) {
            return switch (direction) {
                case UP -> this.up();
                case DOWN -> this.down();
                case LEFT -> this.left();
                case RIGHT -> this.right();
                default -> null;
            };
        }

        public Boolean isTouching(Position other) {
            return Math.abs(this.row() - other.row()) <= 1 && Math.abs(this.col() - other.col()) <= 1;
        }

        public Boolean isDiagonalOf(Position other) {
            return this.row().longValue() != other.row().longValue()
                    && this.col().longValue() != other.col().longValue();
        }

        public Direction verticalDirectionOf(Position other) {
            if (this.col() < other.col()) {
                return Direction.LEFT;
            }

            if (this.col() > other.col()) {
                return Direction.RIGHT;
            }

            return null;
        }

        public Direction horizontalDirectionOf(Position other) {
            if (this.row() < other.row()) {
                return Direction.UP;
            }

            if (this.row() > other.row()) {
                return Direction.DOWN;
            }

            return null;
        }

        public Direction reverseHorizontalDirectionOf(Position other) {
            if (this.row() < other.row()) {
                return Direction.DOWN;
            }

            if (this.row() > other.row()) {
                return Direction.UP;
            }

            return null;
        }

        public Direction reverseVerticalDirectionOf(Position other) {
            if (this.col() < other.col()) {
                return Direction.RIGHT;
            }

            if (this.col() > other.col()) {
                return Direction.LEFT;
            }

            return null;
        }

        public Position up() {
            return new Position(row - 1, col);
        }

        public Position down() {
            return new Position(row + 1, col);
        }

        public Position right() {
            return new Position(row, col + 1);
        }

        public Position left() {
            return new Position(row, col - 1);
        }
    }

    private static class Knot {
        private final Map<Position, Rope> tailMoves;
        private Rope head;
        private Rope tail;
        private Knot next;
        private Map<Position, Rope> moves;

        public Knot(Rope head, Rope tail, Knot next, Map<Position, Rope> moves, Map<Position, Rope> tailMoves) {
            this.head = head;
            this.tail = tail;
            this.next = next;
            this.moves = moves;
            this.tailMoves = tailMoves;
        }

        public void move(Direction direction, Long steps) {
            var stepsLeft = steps;

            while (stepsLeft != 0) {
                stepsLeft--;
                move(direction);
            }
        }

        public void move(Direction direction) {
            final var nextPosition = head.position().in(direction);
            if (!moves.containsKey(nextPosition)) {
                final var nextMove = new Rope(nextPosition, new HashMap<>());
                moves.put(nextPosition, nextMove);
            }

            if (moves.containsKey(nextPosition)) {
                head.next().put(direction, moves.get(nextPosition));
            }

            if (head.next().containsKey(direction)) {
                head = head.next().get(direction);
            }

            if (!head.position().isTouching(tail.position())) {
                computeTailMovement(direction);
            }
        }

        public Direction computeTailMovement(Direction direction) {
            var nextPosition = tail.position().in(direction);

            if (head.position().isDiagonalOf(tail.position())) {
                final var secondaryDirection = findSecondaryDirection(direction);
                nextPosition = nextPosition.in(secondaryDirection);
                final var diagonalDirection = Direction.of(direction, secondaryDirection);
                moveTail(diagonalDirection, nextPosition);
                return secondaryDirection;
            } else {
                moveTail(direction, nextPosition);
            }

            return direction;

        }

        private void moveTail(Direction direction, Position nextPosition) {
            if (!tailMoves.containsKey(nextPosition)) {
                final var nextMove = new Rope(nextPosition, new HashMap<>());
                tailMoves.put(nextPosition, nextMove);
            }

            if (tailMoves.containsKey(nextPosition)) {
                tail.next().put(direction, tailMoves.get(nextPosition));
            }

            if (tail.next().containsKey(direction)) {
                tail = tail.next().get(direction);
            }
        }

        public Direction computeDirection(Direction direction) {
            if (head.position().isDiagonalOf(tail.position())) {
                return direction;
            }
            final var horizontalDirection = head.position().horizontalDirectionOf(tail.position());
            if (horizontalDirection != null) {
                return horizontalDirection;
            }
            final var verticalDirection = head.position().verticalDirectionOf(tail.position());
            if (verticalDirection == null) {
                System.out.println("lol");
            }
            return verticalDirection;

        }

        public Boolean isHeadTouchingTail() {
            return head.position().isTouching(tail.position());
        }

        private Direction findSecondaryDirection(Direction direction) {
            return switch (direction) {
                case UP -> head.position().verticalDirectionOf(tail.position());
                case DOWN -> head.position().verticalDirectionOf(tail.position());
                case LEFT -> head.position().horizontalDirectionOf(tail.position());
                case RIGHT -> head.position().horizontalDirectionOf(tail.position());
                default -> null;
            };
        }

        public Rope head() {
            return head;
        }

        public Rope tail() {
            return tail;
        }

        public Knot next() {
            return next;
        }

        public Map<Position, Rope> moves() {
            return moves;
        }

        public Map<Position, Rope> tailMoves() {
            return tailMoves;
        }
    }

    private class Bridge {
        private final Boolean visualize;

        private Knot head;
        private Knot tail;

        public Bridge() {
            init(1);
            this.visualize = false;
        }

        public Bridge(Boolean visualize) {
            init(1);
            this.visualize = visualize;
        }

        public Bridge(Integer knots) {
            init(knots);
            this.visualize = false;
        }

        public Bridge(Integer knots, Boolean visualize) {
            init(knots);
            this.visualize = visualize;
        }

        private Knot initializeKnots(Knot head, Integer knots) {
            var knot = head;
            var knotsCount = knots - 1;
            while (knotsCount != 0) {
                knotsCount--;
                knot.next = setupKnot(knot);
                knot = knot.next;
            }

            return knot;
        }

        private void init(Integer knots) {
            this.head = setupKnot();
            this.tail = initializeKnots(this.head, knots);
        }

        private Knot setupKnot() {
            final var startingPosition = new Position(0L, 0L);
            final var start = new Rope(startingPosition, new HashMap<>());
            return new Knot(start, start, null, new HashMap<>(Map.of(startingPosition, start)),
                    new HashMap<>(Map.of(startingPosition, start)));
        }

        private Knot setupKnot(Knot prev) {
            final var startingPosition = new Position(0L, 0L);
            final var start = new Rope(startingPosition, new HashMap<>());
            return new Knot(prev.tail(), start, null, new HashMap<>(Map.of(startingPosition, start)),
                    new HashMap<>(Map.of(startingPosition, start)));
        }

        public void move(Direction direction, Long steps) {
            var stepsLeft = steps;

            while (stepsLeft != 0) {
                visualize();
                stepsLeft--;

                head.move(direction);
                moveLeadingKnots(head, direction);
            }
        }

        private void moveLeadingKnots(Knot start, Direction direction) {
            var knot = start;
            var previousDirection = direction;

            while (knot.next != null) {
                knot = updateAndGetNext(knot);

                if (knot.isHeadTouchingTail()) {
                    break;
                }

                final var computedDirection = knot.computeDirection(previousDirection);
                previousDirection = knot.computeTailMovement(computedDirection);
            }
        }

        private Knot updateAndGetNext(Knot knot) {
            knot.next.head = knot.tail;
            knot.next.moves = knot.tailMoves;
            return knot.next;
        }

        private void visualize() {
            if (visualize) {
                BridgePrinter.printCurrentState(this);
            }
        }
    }

    private class BridgePrinter {

        private static final Position start = new Position(0L, 0L);

        private static final String emptyMark = ".";
        private static final String startMark = "s";
        private static final String moveMark = "#";
        private static final String headMark = "H";
        private static final String tailMark = "T";

        private static final Long fromRow = -14L;
        private static final Long toRow = 6L;
        private static final Long fromCol = -11L;
        private static final Long toCol = 15L;

        // private static final Long fromRow = -50L;
        // private static final Long toRow = 50L;
        // private static final Long fromCol = -100L;
        // private static final Long toCol = 100L;

        public static void printCurrentState(Bridge bridge) {
            final var map = createMap(fromRow, toRow, fromCol, toCol);

            markStart(map);
            markTailMoves(bridge, map);
            markKnots(bridge.head, map);

            printMap(map);
            System.out.println();
        }

        private static void markKnots(Knot start, Map<Long, Map<Long, String>> map) {
            var knot = start;
            var level = 0;

            while (knot != null) {
                drawKnot(knot.head, level, map);

                level++;

                if (knot.next == null) {
                    drawTail(knot.tail, level, map);
                }

                knot = knot.next;
            }
        }

        private static void drawKnot(Rope rope, Integer level, Map<Long, Map<Long, String>> map) {
            draw(rope.position(), map, (prev) -> uniqueMarking(level, prev));
        }

        private static void drawTail(Rope rope, Integer level, Map<Long, Map<Long, String>> map) {
            draw(rope.position(), map, (prev) -> priorityMarking(level, prev, (num) -> tailMark));
        }

        private static String uniqueMarking(Integer knot, String prev) {
            return priorityMarking(knot, prev, BridgePrinter::uniqueMark);
        }

        private static String priorityMarking(Integer knot, String prev, Function<Integer, String> uniqueMark) {
            return switch (prev) {
                case emptyMark, moveMark, startMark -> uniqueMark.apply(knot);
                default -> prev;
            };
        }

        public static String uniqueMark(Integer knot) {
            if (knot == 0) {
                return headMark;
            }

            return String.valueOf(knot);
        }

        public static void printFinalState(Bridge bridge) {
            final var map = createMap(fromRow, toRow, fromCol, toCol);
            markTailMoves(bridge, map);
            markStart(map);
            printMap(map);
        }

        public static void markStart(Map<Long, Map<Long, String>> map) {
            draw(start, map, (prev) -> startMark);
        }

        public static void markTailMoves(Bridge bridge, Map<Long, Map<Long, String>> map) {
            bridge.tail.tailMoves.values()
                    .forEach(move -> draw(move.position(), map, (prev) -> moveMark));
        }

        private static void draw(Position position, Map<Long, Map<Long, String>> map, Function<String, String> newValue) {
            map.computeIfPresent(position.row(), (row, columns) -> {
                columns.computeIfPresent(position.col(), (col, column) -> newValue.apply(column));
                return columns;
            });
        }

        public static void printMap(Map<Long, Map<Long, String>> map) {
            final var stringBuilder = new StringBuilder();

            map.values().forEach(row -> {
                row.values().forEach(stringBuilder::append);
                stringBuilder.append("\n");
            });

            System.out.println(stringBuilder);
            sleep();
        }

        private static void sleep() {
            try {
                Thread.sleep(Math.round(1000 * frequency));
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }

        public static Map<Long, Map<Long, String>> createMap(Long fromRow, Long toRow, Long fromCol, Long toCol) {
            final var map = new TreeMap<Long, Map<Long, String>>();
            LongStream.range(fromRow, toRow)
                    .forEach(row -> map.putIfAbsent(row, createColumnMap()));

            return map;
        }

        private static Map<Long, String> createColumnMap() {
            final var columns = LongStream.range(fromCol, toCol).boxed()
                    .collect(Collectors.toMap(col -> col, col -> emptyMark));

            return new TreeMap<Long, String>(columns);
        }
    }
}