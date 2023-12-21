package solutions.year2022;

import common.Problem;
import common.SolutionRunner;
import common.SplitSolution;

import java.util.LinkedList;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.stream.Collectors;
import java.util.stream.IntStream;
import java.util.stream.Stream;

public class Day20 extends Problem {

    public static final String file = "2022/20.txt";
    public static final Problem INSTANCE = new Day20(file);
    public static final SolutionRunner solutionRunner = new SolutionRunner(INSTANCE);

    public Day20(String file) {
        super(file);
    }

    public Day20(String file, SplitSolution.SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    public static void main(String[] args) throws Exception {
        solutionRunner.run();
    }

    @Override
    public String solutionPartOne(Stream<String> data) {
// final var reader = new BufferedReader(new StringReader("""
        // 1
        // 2
        // -3
        // 3
        // -2
        // 0
        // 4"""));
        final var counter = new AtomicInteger(0);

        final LinkedList<NumberDone> numbers = data
                .map(num -> new NumberDone(Integer.parseInt(num), counter.getAndIncrement()))
                .collect(Collectors.toCollection(LinkedList::new));
        final LinkedList<NumberDone> clonedNumbers = new LinkedList<>(numbers);

        clonedNumbers.forEach(number -> {
            if (number.number() == 0) {
                return;
            }

            var movement = (number.number() + numbers.indexOf(number));

            if (number.number() < 0) {
                movement--;
            }

            if (movement >= numbers.size()) {
                movement %= (numbers.size() - 1);
            }

            if (movement <= -numbers.size()) {
                movement %= numbers.size();
            }

            if (movement < 0) {
                movement += numbers.size();
            }

            numbers.remove(number);
            numbers.add(movement, number);
            // System.out.println(numbers.stream().map(s -> s.number()).toList());
        });

        var deMapped = numbers.stream().map(s -> s.number()).toList().indexOf(0);
        return String.valueOf(IntStream.of(1000, 2000, 3000)
                .map(s -> s + (deMapped % numbers.size()))
                .map(s -> numbers.get(s).number()).sum());
        // System.out.println(numbers.stream().map(s -> s.number()).toList());

        // return String.valueOf(
        // numbers.get(1000).number() + numbers.get(2000).number() +
        // numbers.get(3000).number());

        // numbers.get(3000 - 1).number();
        // return "%s %s %s".formatted(numbers.get(1000 - 1).number(), numbers.get(2000
        // - 1).number(),
        // numbers.get(3000 - 1).number());
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        return null;
    }

    public record NumberDone(Integer number, Integer index) {
    }

}
