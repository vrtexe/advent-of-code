package solutions.year2022;

import common.Problem;

import java.util.ArrayList;
import java.util.List;
import java.util.Arrays;
import java.util.stream.Stream;

public class Day17 extends Problem {

    public static final String file = "2022/17-test.txt";
    public static final Problem INSTANCE = new Day17(file);

    // ####

    enum RockShape {
        I(new String[][] {
                { ".", ".", "@", "@", "@", "@", "." }
        }),
        P(new String[][] {
                { ".", ".", ".", "@", ".", ".", "." },
                { ".", ".", "@", "@", "@", ".", "." },
                { ".", ".", ".", "@", ".", ".", "." }
        }),
        J(new String[][] {
                { ".", ".", ".", ".", "@", ".", "." },
                { ".", ".", ".", ".", "@", ".", "." },
                { ".", ".", "@", "@", "@", ".", "." },
        }),
        L(new String[][] {
                { ".", ".", "@", ".", ".", ".", "." },
                { ".", ".", "@", ".", ".", ".", "." },
                { ".", ".", "@", ".", ".", ".", "." },
                { ".", ".", "@", ".", ".", ".", "." },
        }),
        O(new String[][] {
                { ".", ".", "@", "@", ".", ".", "." },
                { ".", ".", "@", "@", ".", ".", "." },
        });

        private final String[][] data;

        public String[][] data() {
            return data;
        }

        RockShape(String[][] data) {
            this.data = data;
        }
    }

    public Day17(String file) {
        super(file);
    }

    public Day17(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    @Override
    public String solutionPartOne(Stream<String> data) {
        final var inputs = data.flatMap(s -> Arrays.stream(s.split(""))).toList();

        var rocks = new RockShape[] { RockShape.I, RockShape.P, RockShape.J, RockShape.L, RockShape.O };
        var emptyLine = new String[] { ".", ".", ".", ".", ".", ".", "." };

        List<List<String>> state = new ArrayList<>();

        var input = 0;
        var rock = 0;
        while (true) {
            var currentRock = rocks[rock];

            List<List<String>> currentState = new ArrayList<>();

            currentState.addAll(Arrays.stream(currentRock.data).map(s -> Arrays.asList(s)).toList());
            currentState.addAll(List.of(Arrays.asList(emptyLine), Arrays.asList(emptyLine), Arrays.asList(emptyLine)));
            currentState.addAll(state);

            // while (true) {
            //     inputs.get(input);

            //     input = (input + 1) % inputs.size();
            // }

            // rock = (rock + 1) % rocks.length;

            // System.err.println(input);
            // break;

        }

        // return null;
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        return null;
    }

    public static void main(String[] args) throws Exception {
        System.out.println(System.getProperty("platform"));
        INSTANCE.run();
    }
}
