package solutions.year2022;

import java.util.stream.Stream;
import common.Problem;

@SuppressWarnings("unused")
public class Day04 extends Problem {

  private static final String file = "2022/04.txt";
  private static final Problem INSTANCE = new Day04(file, SolutionMethod.PART_TWO);

  public Day04(String file) {
    super(file);
  }

  public Day04(String file, SolutionMethod solutionMethod) {
    super(file, solutionMethod);
  }

  @Override
  public String solutionPartOne(Stream<String> data) {

    final long count = data.map(AssignmentPair::of)
        .filter(AssignmentPair::hasOverlap)
        .count();

    return String.valueOf(count);
  }

  @Override
  public String solutionPartTwo(Stream<String> data) {
    final long count = data.map(AssignmentPair::of)
        .filter(AssignmentPair::hasAnyOverlap)
        .count();

    return String.valueOf(count);
  }

  public AssignmentPair parseCleanupPair(String cleanupPair) {
    return AssignmentPair.of(cleanupPair);
  }

  private record CleanupRange(Long from, Long to) {

    public static CleanupRange of(String cleaningRange) {
      final var range = cleaningRange.split("-");
      return new CleanupRange(Long.valueOf(range[0]), Long.valueOf(range[1]));
    }

    public Boolean areOverlapping(CleanupRange other) {
      return this.hasOverlap(other) || other.hasOverlap(this);
    }

    private Boolean hasOverlap(CleanupRange other) {
      return other.from() <= this.from() && other.to() >= this.to();
    }

    public Boolean areOverlappingAtAll(CleanupRange other) {
      return this.hasAnyOverlap(other) || other.hasAnyOverlap(this);
    }

    private Boolean hasAnyOverlap(CleanupRange other) {
      return this.from() <= other.to() && other.from() <= this.to();
    }
  }

  private record AssignmentPair(CleanupRange left, CleanupRange right) {

    public static AssignmentPair of(String cleanupPair) {
      final var cleanupPairs = cleanupPair.split(",");

      final var left = CleanupRange.of(cleanupPairs[0]);
      final var right = CleanupRange.of(cleanupPairs[1]);

      return new AssignmentPair(left, right);
    }

    public Boolean hasOverlap() {
      return left.areOverlapping(right);
    }

    public Boolean hasAnyOverlap() {
      return left.areOverlappingAtAll(right);
    }
  }

  public static void main(String[] args) throws Exception {
    INSTANCE.run();
  }

}
