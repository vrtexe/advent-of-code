package common;

import java.util.stream.Stream;

public interface SplitSolution {

    String solutionPartOne(Stream<String> data);

    String solutionPartTwo(Stream<String> data);

    default String runBothSolutions(Stream<String> data) {
        final var collectedList = data.toList();

        return """
                Solution 1: %s
                Solution 2: %s
                """.formatted(solutionPartOne(collectedList.stream()),
                solutionPartTwo(collectedList.stream()));
    }

    default SolutionMethod getSolutionMethod() {
        return SolutionMethod.ALL;
    }

    enum SolutionMethod {
        ALL, PART_ONE, PART_TWO
    }
}
