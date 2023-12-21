package solutions.year2022;

import common.Problem;

import java.util.*;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;
import java.util.stream.LongStream;
import java.util.stream.Stream;

@SuppressWarnings("unused")
public class Day05 extends Problem {

    public static final String file = "2022/05.txt";
    public static final Problem INSTANCE = new Day05(file);

    public Day05(String file) {
        super(file);
    }

    public Day05(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    public String solutionPartOne(Stream<String> data) {
        final var inputs = data.collect(Collectors.joining("\n")).split("\\n\\n");
        final var supplyStacks = parseCurrentState(inputs[0]);

        return moveCargoByInstructions(supplyStacks, Arrays.asList(inputs[1].split("\n"))).values().stream()
                .map(Stack::peek).collect(Collectors.joining(""));
    }

    public String solutionPartTwo(Stream<String> data) {
        final var inputs = data.collect(Collectors.joining("\n")).split("\\n\\n");
        final var supplyStacks = parseCurrentState(inputs[0]);

        return moveMultiCargoByInstructions(supplyStacks, Arrays.asList(inputs[1].split("\n"))).values().stream()
                .map(Stack::peek).collect(Collectors.joining(""));
    }

    private Map<String, Stack<String>> parseCurrentState(String currentState) {

        final var stacks = Arrays.asList(currentState.split("\n"));
        Collections.reverse(stacks);

        final var stackItems = stacks.subList(1, stacks.size());

        final var result = Arrays.stream(stacks.get(0).split(" +")).skip(1)
                .collect(Collectors.toMap(s -> s, s -> new Stack<String>()));

        stackItems.stream().map(items -> items.split("(?<=\\G.{3} )"))
                .map(Arrays::asList)
                .forEach(items -> {
                    final var stackPosition = new AtomicInteger();

                    items.forEach(item -> {
                        final var resultStack = String.valueOf(stackPosition.incrementAndGet());

                        if (result.containsKey(resultStack) && !item.isBlank()) {
                            result.get(resultStack).push(String.valueOf(item.charAt(1)));
                        }
                    });
                });

        return result;
    }

    private Map<String, Stack<String>> moveCargoByInstructions(Map<String, Stack<String>> cargo,
                                                               List<String> instructions) {
        Pattern pattern = Pattern.compile("move ([0-9]+) from ([0-9]+) to ([0-9]+)");

        instructions.forEach(instruction -> {
            Matcher matcher = pattern.matcher(instruction);
            if (matcher.find()) {
                final var count = Long.parseLong(matcher.group(1));
                final var from = cargo.get(matcher.group(2));
                final var to = cargo.get(matcher.group(3));

                LongStream.range(0, count)
                        .forEach(__ -> to.push(from.pop()));
            }

        });

        return cargo;
    }

    private Map<String, Stack<String>> moveMultiCargoByInstructions(Map<String, Stack<String>> cargo,
                                                                    List<String> instructions) {
        Pattern pattern = Pattern.compile("move ([0-9]+) from ([0-9]+) to ([0-9]+)");

        instructions.forEach(instruction -> {
            Matcher matcher = pattern.matcher(instruction);
            if (matcher.find()) {
                final var count = Integer.parseInt(matcher.group(1));
                final var from = cargo.get(matcher.group(2));
                final var to = cargo.get(matcher.group(3));

                final var crates = LongStream.range(0, count)
                        .mapToObj(__ -> from.pop()).collect(Collectors.toCollection(Stack::new));

                LongStream.range(0, count).forEach(s -> to.push(crates.pop()));

            }

        });

        return cargo;
    }
}
