package common;


import java.util.regex.Pattern;
import java.util.stream.Stream;

public abstract class Problem implements Solution {

    public final String file;
    public final SolutionRunner solutionRunner;
    public final SolutionMethod solutionMethod;

    public Problem() {
        this.file = getDefaultFileName();
        this.solutionMethod = SolutionMethod.ALL;
        this.solutionRunner = new SolutionRunner(this);
    }

    public Problem(SolutionMethod solutionMethod) {
        this.file = getDefaultFileName();
        this.solutionMethod = solutionMethod;
        this.solutionRunner = new SolutionRunner(this);
    }

    public Problem(String file) {
        this.file = file;
        this.solutionMethod = SolutionMethod.ALL;
        this.solutionRunner = new SolutionRunner(this);

    }

    public Problem(String file, SolutionMethod solutionMethod) {
        this.file = file;
        this.solutionMethod = solutionMethod;
        this.solutionRunner = new SolutionRunner(this);
    }

    public void run() throws Exception {
        solutionRunner.run();
    }

    @Override
    public String getFileName() {
        return file;
    }

    @Override
    public SolutionMethod getSolutionMethod() {
        return solutionMethod;
    }

    abstract public String solutionPartOne(Stream<String> data);

    abstract public String solutionPartTwo(Stream<String> data);

    private String getDefaultFileName() {
        final var year = getYearFromPackage();
        final var day = getDayFromClassName();

        return "%s/%s.txt".formatted(year, day);
    }

    private String getYearFromPackage() {
        final var pattern = Pattern.compile("([0-9]+)");
        final var matcher = pattern.matcher(this.getClass().getPackageName());

        if (matcher.find()) {
            return matcher.group(1);
        }

        return "";
    }

    private String getDayFromClassName() {
        final var pattern = Pattern.compile("([0-9]+)");
        final var matcher = pattern.matcher(this.getClass().getSimpleName());

        if (matcher.find()) {
            return matcher.group(1);
        }

        return "";
    }
}
