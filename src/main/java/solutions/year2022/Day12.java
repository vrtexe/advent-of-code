package solutions.year2022;

import common.Problem;

import java.util.*;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.Function;
import java.util.function.Supplier;
import java.util.stream.Stream;

@SuppressWarnings("unused")
public class Day12 extends Problem {

    public static final String file = "2022/12.txt";
    public static final Problem INSTANCE = new Day12(file);

    public Day12(String file) {
        super(file);
    }

    public Day12(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    private static Stream<String> split(String line, @SuppressWarnings("SameParameterValue") String delim) {
        return Arrays.stream(line.split(delim)).map(String::trim);
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    @Override
    public String solutionPartOne(Stream<String> data) {
        final var heightMap = createMap(data);

        return String.valueOf(heightMap.getShortestPath(heightMap.getStart(), heightMap.getEnd()));
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        final var heightMap = createMap(data);

        final var lengths = heightMap.getShortestPathReverse(heightMap.getEnd());
        final var minLength = lengths.entrySet().stream()
                .filter(s -> s.getKey().elevation().equals('a'))
                .mapToLong(Map.Entry::getValue)
                .min()
                .orElse(-1L);

        return String.valueOf(minLength);
    }

    private HeightMap createMap(Stream<String> data) {
        final var map = new HashMap<Position, Square>();
        final var rowNumber = new AtomicLong();

        final var heightMapBuilder = new HeightMapBuilder();

        data.map(line -> split(line, "").map(s -> s.charAt(0)))
                .forEach(row -> {
                    final var colNumber = new AtomicLong();
                    row.forEach(col -> {
                        final var position = new Position(rowNumber.get(), colNumber.getAndIncrement());

                        heightMapBuilder.start(col, position).end(col, position);

                        map.putIfAbsent(position, Square.of(col, position, map));
                    });
                    rowNumber.incrementAndGet();
                });

        heightMapBuilder.squares(map);

        return heightMapBuilder.createHeightMap();
    }

    private enum Direction {
        UP,
        DOWN,
        LEFT,
        RIGHT
    }

    private record HeightMap(Position start, Position end, Map<Position, Square> squares) {

        public Square getStart() {
            return squares.get(start);
        }

        public Square getEnd() {
            return squares.get(end);
        }

        public Long getShortestPath(Square start, Square end) {
            final var result = getShortestPath(start, Square::getNeighbors);
            if (result.containsKey(end)) {
                return result.get(end);
            }
            return Long.MAX_VALUE;
        }

        public Map<Square, Long> getShortestPathReverse(Square start) {
            return getShortestPath(start, Square::getReverseNeighbors);
        }

        public Map<Square, Long> getShortestPath(Square start) {
            return getShortestPath(start, Square::getNeighbors);
        }

        public Map<Square, Long> getShortestPath(Square start, Function<Square, List<Square>> getNeighbors) {
            final var distances = new HashMap<Square, Long>();
            final var visited = new HashSet<Square>();
            final var remaining = new PriorityQueue<Square>(Comparator.comparing(distances::get));

            @SuppressWarnings("MismatchedQueryAndUpdateOfCollection") final var previous = new HashMap<Square, Square>();

            distances.put(start, 0L);
            visited.add(start);
            remaining.offer(start);

            while (!remaining.isEmpty()) {
                final var square = remaining.poll();
                getNeighbors.apply(square).forEach(neighbor -> {
                    final var newPath = distances.get(square) + 1;
                    final var oldPath = distances.get(neighbor);

                    if (!distances.containsKey(neighbor) || newPath < oldPath) {
                        distances.put(neighbor, newPath);
                        previous.put(neighbor, square);
                    }

                    if (!visited.contains(neighbor)) {
                        remaining.offer(neighbor);
                        visited.add(neighbor);
                    }
                });
            }

            return distances;
        }
    }

    private record Square(Character elevation, Position position, Neighbors neighbors) {
        public static Square of(Character character, Position position, Map<Position, Square> map) {
            final var elevation = switch (character) {
                case 'S' -> 'a';
                case 'E' -> 'z';
                default -> character;
            };

            return new Square(elevation, position, Neighbors.of(position, map));
        }

        public List<Square> getNeighbors() {
            return neighbors.get().stream()
                    .filter(this::isNeighbor)
                    .toList();
        }

        public List<Square> getReverseNeighbors() {
            return neighbors.get().stream()
                    .filter(this::isReverseNeighbor)
                    .toList();
        }

        private Boolean isNeighbor(Square neighbor) {
            return neighbor.elevation() - this.elevation() <= 1;
        }

        private Boolean isReverseNeighbor(Square neighbor) {
            return this.elevation() - neighbor.elevation() <= 1;
        }
    }

    private record Neighbors(Supplier<Square> up, Supplier<Square> down, Supplier<Square> left,
                             Supplier<Square> right) {

        public static Neighbors of(Position position, Map<Position, Square> map) {
            return new Neighbors(
                    () -> map.get(position.up()),
                    () -> map.get(position.down()),
                    () -> map.get(position.left()),
                    () -> map.get(position.right()));
        }

        public List<Square> get() {
            return Stream.of(up.get(), down.get(), left.get(), right.get())
                    .filter(Objects::nonNull)
                    .toList();
        }

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

    @SuppressWarnings("UnusedReturnValue")
    private static final class HeightMapBuilder {
        Position start;
        Position end;
        Map<Position, Square> squares;

        public HeightMap createHeightMap() {
            return new HeightMap(start, end, squares);
        }

        public HeightMapBuilder start(Character elevation, Position start) {
            if (elevation == 'S') {
                this.start = start;
            }

            return this;
        }

        public HeightMapBuilder end(Character elevation, Position end) {
            if (elevation == 'E') {
                this.end = end;
            }
            return this;
        }

        public HeightMapBuilder squares(Map<Position, Square> squares) {
            this.squares = squares;
            return this;
        }

    }
}
