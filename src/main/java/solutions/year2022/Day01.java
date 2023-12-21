package solutions.year2022;

import common.Problem;

import java.util.Arrays;
import java.util.stream.Collectors;
import java.util.stream.Stream;

public class Day01 extends Problem {

    public static final String file = "2022/01.txt";
    public static final Problem INSTANCE = new Day01(file);

    public Day01(String file) {
        super(file);
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    @Override
    public String solutionPartOne(Stream<String> data) {
        var maxCalories = Arrays.stream(data.collect(Collectors.joining("\n")).split("\\n\\n"))
                .map(elf -> Arrays.asList(elf.split("\n")))
                .mapToLong(elf -> elf.stream().mapToLong(Long::parseLong).sum())
                .max()
                .orElse(0L);

        return String.valueOf(maxCalories);
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        var maxCalories = Arrays.stream(data.collect(Collectors.joining("\n")).split("\\n\\n"))
                .map(elf -> Arrays.asList(elf.split("\n")))
                .map(elf -> elf.stream().mapToLong(Long::parseLong).sum())
                .sorted((caloriesElf, caloriesElfOther) -> Long.compare(caloriesElfOther, caloriesElf))
                .mapToLong(calories -> calories).limit(3).sum();

        return String.valueOf(maxCalories);
    }
}
