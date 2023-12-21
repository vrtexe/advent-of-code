package common;

import java.util.stream.Stream;

public interface Solution extends SplitSolution {

    default String solve(Stream<String> data) {
        return switch (this.getSolutionMethod()) {
            case ALL -> runBothSolutions(data);
            case PART_ONE -> solutionPartOne(data);
            case PART_TWO -> solutionPartTwo(data);
        };
    }

    String getFileName();
}
