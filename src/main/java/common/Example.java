package common;

import java.util.stream.Stream;


public class Example extends Problem {

    public static final String file = "src/resources/01.txt";
    public static final Problem INSTANCE = new Example(file);

    public Example(String file) {
        super(file);
    }

    public Example(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    @Override
    public String solutionPartOne(Stream<String> data) {
        return null;
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        return null;
    }
}
