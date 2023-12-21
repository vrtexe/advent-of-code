package solutions.year2022;

import common.Problem;

import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.IntStream;
import java.util.stream.LongStream;
import java.util.stream.Stream;

public class Day14 extends Problem {

    public static final String file = "2022/14.txt";
    public static final Problem INSTANCE = new Day14(file, SolutionMethod.ALL);

    private final Position sandSource = new Position(500L, 0L);

    // private long startX = 493 - 50;
    // private long endX = 504 + 50;

    // private long startY = 0;
    // private long endY = 10 + 4;

    private final long startX = 493 - 200;
    private final long endX = 504 + 200;

    private final long startY = 0;
    private final long endY = 10 + 160;

    private final long leftLimit = 200;
    private final long rightLimit = 200;

    private final double frequency = 0;
    private final boolean visual = false;
    private final boolean heading = false;

    private Long lastY = 0L;

    public Day14(String file) {
        super(file);
    }

    public Day14(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    @Override
    public String solutionPartOne(Stream<String> data) {
        final var waterfall = mapRockPositions(data);
        lastY = findLastY(waterfall);

        waterfall.put(sandSource, Type.SOURCE);

        dropAllPossibleSands(waterfall);

        final var restingSand = waterfall.values().stream().filter(s -> s == Type.SAND).count();

        return String.valueOf(restingSand);
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        final var waterfall = mapRockPositions(data);
        mapFloorPositions(waterfall);
        lastY = findLastY(waterfall);

        waterfall.put(sandSource, Type.SOURCE);

        dropAllPossibleSands(waterfall);

        final var restingSand = waterfall.values().stream().filter(s -> s == Type.SAND).count();

        return String.valueOf(restingSand);
    }

    private void dropAllPossibleSands(Map<Position, Type> waterfall) {
        var shouldContinue = true;

        while (shouldContinue) {
            shouldContinue = generateSandAndMoveDown(waterfall);
        }
    }

    private Boolean generateSandAndMoveDown(Map<Position, Type> waterfall) {
        var sandUnit = sandSource;

        while (sandCanFall(waterfall, sandUnit)) {
            printPositions(waterfall);
            sandUnit = fall(waterfall, sandUnit);

            if (sandUnit.y() > lastY + 4) {
                waterfall.remove(sandUnit);
                return false;
            }
        }

        if (sandUnit.equals(sandSource)) {
            waterfall.put(sandUnit, Type.SAND);
            printPositions(waterfall);
            return false;
        }

        printPositions(waterfall);

        return true;
    }

    private Long findLastY(Map<Position, Type> waterfall) {
        return waterfall.entrySet()
                .stream().filter(s -> s.getValue() == Type.ROCK)
                .mapToLong(s -> s.getKey().y())
                .max()
                .orElse(0L);
    }

    private Long findFirstY(Map<Position, Type> waterfall) {
        return waterfall.entrySet()
                .stream().filter(s -> s.getValue() == Type.ROCK)
                .mapToLong(s -> s.getKey().y())
                .min()
                .orElse(0L);
    }

    private Long findFirstX(Map<Position, Type> waterfall) {
        return waterfall.entrySet()
                .stream().filter(s -> s.getValue() == Type.ROCK)
                .mapToLong(s -> s.getKey().x())
                .min()
                .orElse(0L);
    }

    private Long findLastX(Map<Position, Type> waterfall) {
        return waterfall.entrySet()
                .stream().filter(s -> s.getValue() == Type.ROCK)
                .mapToLong(s -> s.getKey().x())
                .max()
                .orElse(0L);
    }

    private Position fall(Map<Position, Type> waterfall, Position sandUnit) {
        if (waterfall.get(sandUnit) == Type.SAND) {
            waterfall.remove(sandUnit);
        }

        final var down = sandUnit.down();
        if (!waterfall.containsKey(down)) {
            waterfall.put(down, Type.SAND);
            return down;
        }

        final var downLeft = down.left();
        if (!waterfall.containsKey(downLeft)) {
            waterfall.put(downLeft, Type.SAND);
            return downLeft;
        }

        final var downRight = down.right();
        if (!waterfall.containsKey(downRight)) {
            waterfall.put(downRight, Type.SAND);
            return downRight;
        }

        return sandUnit;
    }

    private Boolean sandCanFall(Map<Position, Type> waterfall, Position sandUnit) {
        final var down = sandUnit.down();
        return !(waterfall.containsKey(down) && waterfall.containsKey(down.left())
                && waterfall.containsKey(down.right()));
    }

    private Map<Position, Type> mapFloorPositions(Map<Position, Type> waterfall) {
        final var floor = findLastY(waterfall) + 2;
        final var from = findFirstX(waterfall) - leftLimit;
        final var to = findLastX(waterfall) + rightLimit;

        LongStream.range(from, to)
                .forEach(x -> waterfall.put(new Position(x, floor), Type.ROCK));

        return waterfall;
    }

    private Map<Position, Type> mapRockPositions(Stream<String> data) {
        return data.flatMap(line -> findPositionsBetween(Arrays.asList(line.split(" -> "))))
                .distinct()
                .collect(Collectors.toMap(position -> position, position -> Type.ROCK));
    }

    private Stream<Position> findPositionsBetween(List<String> paths) {

        final var positions = paths.stream().map(s -> Arrays.asList(s.split(",")))
                .map(this::getPositionFromString).toList();

        return IntStream.range(1, positions.size())
                .mapToObj(index -> getPositionsForPath(positions.get(index - 1), positions.get(index)))
                .flatMap(s -> s.stream())
                .distinct();
    }

    private List<Position> getPositionsForPath(Position from, Position to) {
        final var direction = getPathDirection(from, to);
        final var positions = new ArrayList<Position>(List.of(from, to));

        var position = from.in(direction);
        while (!position.equals(to)) {
            positions.add(position);
            position = position.in(direction);
        }

        return positions;
    }

    private Direction getPathDirection(Position from, Position to) {
        final var direction = getVerticalDirection(from, to);

        if (direction == null) {
            return getHorizontalDirection(from, to);
        }

        return direction;
    }

    private Direction getVerticalDirection(Position from, Position to) {
        final var direction = from.y() - to.y();

        if (direction < 0) {
            return Direction.DOWN;
        }

        if (direction > 0) {
            return Direction.UP;
        }

        return null;
    }

    private Direction getHorizontalDirection(Position from, Position to) {
        final var direction = from.x() - to.x();

        if (direction < 0) {
            return Direction.RIGHT;
        }

        if (direction > 0) {
            return Direction.LEFT;
        }

        return null;
    }

    private Position getPositionFromString(List<String> position) {
        return new Position(Long.parseLong(position.get(0)), Long.parseLong(position.get(1)));
    }

    private Map<Long, Map<Long, String>> mapStringPositions(Map<Position, Type> positions) {
        final var map = new TreeMap<Long, Map<Long, String>>();

        LongStream.range(startY, endY)
                .forEach(y -> {
                    map.put(y, new TreeMap<>());
                    LongStream.range(startX, endX)
                            .forEach(x -> {
                                final var newPosition = new Position(x, y);
                                if (positions.containsKey(newPosition)) {
                                    map.get(y).put(x, positions.get(newPosition).getValue());
                                } else {
                                    map.get(y).put(x, Type.AIR.getValue());
                                }
                            });
                });

        return map;
    }

    private void printPositions(Map<Position, Type> positions) {
        printPositions(positions, false);
    }

    private void printPositions(Map<Position, Type> positions, Boolean override) {
        if (!visual && !override) {
            return;
        }

        final var map = mapStringPositions(positions);
        final var stringBuilder = new StringBuilder();

        map.entrySet().forEach(y -> {

            if (heading) {
                stringBuilder.append(getHeadingString(y));
                stringBuilder.append(y.getKey());
            }

            y.getValue().values().forEach(x -> {
                stringBuilder.append(x);
            });

            stringBuilder.append("\n");
        });

        System.out.println(stringBuilder.toString());
        sleep();
    }

    private String getHeadingString(Map.Entry<Long, Map<Long, String>> y) {
        final var stringBuilder = new StringBuilder();
        if (y.getKey() == 0) {
            final var maxLength = String.valueOf(endY).length();
            IntStream.range(0, maxLength + 1)
                    .mapToObj(num -> y.getValue().keySet().stream().map(s -> tryGetCharOrDefault(s.toString(), num)))
                    .forEach(s -> {
                        // System.out.print(" ");
                        stringBuilder.append(" ");
                        s.forEach(stringBuilder::append);
                        stringBuilder.append("\n");
                        // System.out.println();
                    });
        }
        return stringBuilder.toString();
    }

    private Character tryGetCharOrDefault(String string, Integer position) {
        try {
            return string.charAt(position);
        } catch (Exception e) {
            return ' ';
        }
    }

    private void sleep() {
        try {
            Thread.sleep(Math.round(1000 * frequency));
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }

    private enum Type {
        ROCK("#"),
        AIR("."),
        SOURCE("+"),
        SAND("o");

        private final String value;

        private Type(String value) {
            this.value = value;
        }

        private String getValue() {
            return value;
        }
    }

    private enum Direction {
        UP,
        DOWN,
        LEFT,
        RIGHT
    }

    private record Position(Long x, Long y) {

        public Position in(Direction direction) {
            return switch (direction) {
                case UP -> this.up();
                case DOWN -> this.down();
                case LEFT -> this.left();
                case RIGHT -> this.right();
                default -> null;
            };
        }

        public Position up() {
            return new Position(x, y - 1);
        }

        public Position down() {
            return new Position(x, y + 1);
        }

        public Position right() {
            return new Position(x + 1, y);
        }

        public Position left() {
            return new Position(x - 1, y);
        }

    }
}
