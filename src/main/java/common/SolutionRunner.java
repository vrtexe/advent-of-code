package common;


public class SolutionRunner {

    private final Solution solution;

    public SolutionRunner(Solution solution) {
        this.solution = solution;
    }

    public void run() throws Exception {

        final var file = new FileLoader(solution.getFileName());

        System.out.println(solution.solve(file.loadFile()));

        file.closeStream();
    }
}
