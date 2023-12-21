package solutions.year2022;

import common.Problem;

import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.BiConsumer;
import java.util.function.Function;
import java.util.stream.Collectors;
import java.util.stream.IntStream;
import java.util.stream.Stream;

@SuppressWarnings("unused")
public class Day08 extends Problem {

    public static final String file = "2022/08.txt";
    public static final Problem INSTANCE = new Day08(file, SolutionMethod.ALL);

    private static final int ROWS = 98;
    private static final int COLUMNS = 98;

    public Day08(String file) {
        super(file);
    }

    public Day08(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    @Override
    public String solutionPartOne(Stream<String> data) {
        final var trees = mapTrees(data);

        trees.values().forEach(tree -> tree.trees = getSurroundingTrees(tree, trees));

        Arrays.stream(Side.values()).forEach(side -> checkVisibilityFromSide(trees, side));

        return String.valueOf(trees.values().stream().filter(this::isAnyVisible).count());
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        final var trees = mapTrees(data);

        trees.values().forEach(tree -> tree.trees = getSurroundingTrees(tree, trees));

        final var highestScenicScore = trees.values().stream()
                .filter(Tree::isNotEdge)
                .mapToLong(Tree::getScenicScore)
                .max()
                .orElse(0L);

        return String.valueOf(highestScenicScore);
    }

    private void printValues(Map<Position, Tree> trees) {
        IntStream.range(0, ROWS + 1)
                .mapToObj(row -> trees.get(new Position(row, 0)))
                .forEach(tree -> {
                    var current = tree;
                    while (current.trees().right() != null) {
                        System.out.printf("%s", isAnyVisible(current) ? 1 : 0);
                        current = current.trees().right();
                    }
                    System.out.printf("%s", isAnyVisible(current) ? 1 : 0);
                    System.out.println();
                });
    }

    private void checkVisibilityFromSide(Map<Position, Tree> trees, Side side) {
        switch (side) {
            case TOP -> computeVisibility(
                    IntStream.range(1, COLUMNS).mapToObj(s -> trees.get(new Position(ROWS, s))).toList(),
                    t -> t.trees().top(), (t, s) -> setVisibility(t, s, side),
                    t -> isLast(t, side));
            case BOTTOM -> computeVisibility(
                    IntStream.range(1, COLUMNS).mapToObj(s -> trees.get(new Position(0, s))).toList(),
                    t -> t.trees().bottom(), (t, s) -> setVisibility(t, s, side),
                    t -> isLast(t, side));
            case LEFT -> computeVisibility(
                    IntStream.range(1, ROWS).mapToObj(s -> trees.get(new Position(s, COLUMNS))).toList(),
                    t -> t.trees().left(), (t, s) -> setVisibility(t, s, side),
                    t -> isLast(t, side));
            case RIGHT -> computeVisibility(
                    IntStream.range(1, ROWS).mapToObj(s -> trees.get(new Position(s, 0))).toList(),
                    t -> t.trees().right(), (t, s) -> setVisibility(t, s, side),
                    t -> isLast(t, side));
        }
    }

    private Boolean isLast(Tree tree, Side side) {
        return switch (side) {
            case TOP -> tree.trees().top() == null;
            case BOTTOM -> tree.trees().bottom() == null;
            case LEFT -> tree.trees().left() == null;
            case RIGHT -> tree.trees().right() == null;
        };
    }

    private void setVisibility(Tree tree, Boolean value, Side side) {
        tree.visible = changeVisibility(tree, value, side);
    }

    private TreeVisibilityDirection changeVisibility(Tree tree, Boolean value, Side side) {
        return switch (side) {
            case TOP ->
                new TreeVisibilityDirection(value, tree.visible.bottom(), tree.visible.left(), tree.visible.right());
            case BOTTOM ->
                new TreeVisibilityDirection(tree.visible.top(), value, tree.visible.left(), tree.visible.right());
            case LEFT ->
                new TreeVisibilityDirection(tree.visible.top(), tree.visible.bottom(), value, tree.visible.right());
            case RIGHT ->
                new TreeVisibilityDirection(tree.visible.top(), tree.visible.bottom(), tree.visible.left(), value);
        };
    }

    private void computeVisibility(List<Tree> start,
            Function<Tree, Tree> next,
            BiConsumer<Tree, Boolean> setVisibility,
            Function<Tree, Boolean> isLast) {
        start.forEach(startTree -> {
            var currentTree = startTree;
            var max = startTree.height();

            while (!isLast.apply(currentTree)) {
                currentTree = next.apply(currentTree);
                if (currentTree.height() > max) {
                    setVisibility.accept(currentTree, true);
                    max = currentTree.height();
                }
            }
        });
    }

    private Boolean isAnyVisible(Tree tree) {
        return isVisible(tree.visible.right) || isVisible(tree.visible.left) ||
                isVisible(tree.visible.top) || isVisible(tree.visible.bottom);
    }

    private Boolean isVisible(Boolean visible) {
        return visible != null && visible;
    }

    private SurroundingTrees getSurroundingTrees(Tree tree, Map<Position, Tree> trees) {
        return new SurroundingTrees(trees.get(tree.position().above()),
                trees.get(tree.position().bellow()),
                trees.get(tree.position().before()),
                trees.get(tree.position().after()));
    }

    private Map<Position, Tree> mapTrees(Stream<String> data) {
        final var rowNumber = new AtomicInteger(-1);

        return data.flatMap(treeRow -> {
            rowNumber.incrementAndGet();
            final var colNumber = new AtomicInteger();
            return Arrays.stream(treeRow.split(""))
                    .map(tree -> switch (rowNumber.get()) {
                        case 0 -> createNewTree(tree, rowNumber.get(), colNumber.getAndIncrement(), Side.TOP);
                        case ROWS -> createNewTree(tree, rowNumber.get(), colNumber.getAndIncrement(), Side.BOTTOM);
                        default -> switch (colNumber.get()) {
                            case 0 -> createNewTree(tree, rowNumber.get(), colNumber.getAndIncrement(), Side.LEFT);
                            case COLUMNS ->
                                createNewTree(tree, rowNumber.get(), colNumber.getAndIncrement(), Side.RIGHT);
                            default -> createNewTree(tree, rowNumber.get(), colNumber.getAndIncrement(), null);
                        };
                    });
        })
                .collect(Collectors.toMap(Tree::position, s -> s));
    }

    private Tree createNewTree(String height, Integer row, Integer col, Side side) {
        final var direction = side == null
                ? new TreeVisibilityDirection(false, false, false, false)
                : switch (side) {
                    case TOP -> new TreeVisibilityDirection(true, false, false, false);
                    case BOTTOM -> new TreeVisibilityDirection(false, true, false, false);
                    case LEFT -> new TreeVisibilityDirection(false, false, true, false);
                    case RIGHT -> new TreeVisibilityDirection(false, false, false, true);
        };

        return new Tree(Long.parseLong(height), new Position(row, col), null, direction);
    }

    private enum Side {
        BOTTOM("BOTTOM"),
        TOP("TOP"),
        LEFT("LEFT"),
        RIGHT("RIGHT");

        final String value;

        Side(String value) {
            this.value = value;
        }

        public String getValue() {
            return value;
        }
    }

    private record Position(Integer row, Integer col) {

        public Position above() {
            return new Position(row - 1, col);
        }

        public Position bellow() {
            return new Position(row + 1, col);
        }

        public Position after() {
            return new Position(row, col + 1);
        }

        public Position before() {
            return new Position(row, col - 1);
        }
    }

    private record SurroundingTrees(Tree top, Tree bottom, Tree left, Tree right) {
    }

    private record TreeVisibilityDirection(Boolean top, Boolean bottom, Boolean left, Boolean right) {
    }

    private static class Tree {
        Long height;
        Position position;
        SurroundingTrees trees;
        TreeVisibilityDirection visible;

        public Tree(Long height, Position position, SurroundingTrees trees, TreeVisibilityDirection visible) {
            this.height = height;
            this.position = position;
            this.trees = trees;
            this.visible = visible;
        }

        public static Function<Tree, Tree> direction(Side side) {
            return switch (side) {
                case TOP -> tree -> tree.trees().top();
                case BOTTOM -> tree -> tree.trees().bottom();
                case LEFT -> tree -> tree.trees().left();
                case RIGHT -> tree -> tree.trees().right();
            };
        }

        public Long height() {
            return height;
        }

        public Position position() {
            return position;
        }

        public SurroundingTrees trees() {
            return trees;
        }

        public TreeVisibilityDirection visible() {
            return visible;
        }

        public Boolean isNotEdge() {
            return this.position().col() != 0
                    && this.position().col() != COLUMNS
                    && this.position().row() != 0
                    && this.position().row() != ROWS;
        }

        public Long getScenicScore() {
            return Arrays.stream(Side.values())
                    .mapToLong(side -> calculateScenicRoute(side, direction(side)))
                    .reduce(1L, (a, b) -> a * b);
        }

        private Long calculateScenicRoute(Side side, Function<Tree, Tree> next) {
            var nextTree = next.apply(this);
            var scenicScore = 1L;
            while (next.apply(nextTree) != null && this.height() > nextTree.height()) {
                scenicScore++;
                nextTree = next.apply(nextTree);
            }

            return scenicScore;
        }
    }
}