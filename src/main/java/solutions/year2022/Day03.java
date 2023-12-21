package solutions.year2022;

import common.Problem;

import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;
import java.util.stream.IntStream;
import java.util.stream.Stream;

@SuppressWarnings("unused")
public class Day03 extends Problem {

    public static final String file = "2022/03.txt";
    public static final Problem INSTANCE = new Day03(file);
    public final Map<String, Integer> priority = Stream
            .concat(Stream.iterate('a', i -> ++i).limit(26), Stream.iterate('A', i -> ++i).limit(26))
            .collect(Collectors.toMap(Object::toString, s -> s - (Character.isLowerCase(s) ? 'a' - 1 : 'A' - 27)));

    public Day03(String file) {
        super(file);
    }

    public Day03(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    @Override
    public String solutionPartOne(Stream<String> data) {
        var result = data.map(this::createRucksack).flatMap(this::getOutliers).mapToInt(priority::get).sum();

        return String.valueOf(result);
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        var rucksacks = data.map(this::createRucksack).toList();

        var result = splitList(rucksacks, 3).stream().flatMap(this::findCommonItem)
                .mapToInt(priority::get).sum();

        return String.valueOf(result);
    }

    private Stream<String> getOutliers(Rucksack rucksack) {
        return Arrays.stream(rucksack.left().split(""))
                .filter(leftItem -> rucksack.right().contains(leftItem))
                .max((a, s) -> priority.get(a).compareTo(priority.get(s))).stream();
    }

    private Rucksack createRucksack(String rucksack) {
        final var halfPoint = rucksack.length() / 2;

        return new Rucksack(rucksack.substring(0, halfPoint), rucksack.substring((halfPoint)),
                rucksack);
    }

    private List<List<Rucksack>> splitList(List<Rucksack> rucksacks, @SuppressWarnings("SameParameterValue") Integer parts) {
        final int lastPart = (rucksacks.size() + parts - 1) / parts;
        return IntStream.range(0, lastPart)
                .mapToObj(part -> rucksacks.subList(part * parts, Math.min(parts * part + parts, rucksacks.size())))
                .toList();
    }

    private Stream<String> findCommonItem(List<Rucksack> rucksacks) {
        final var rucksack = rucksacks.get(0).complete();

        final var as = Arrays.stream(rucksack.split("")).distinct()
                .filter(item -> rucksacks.stream().allMatch(sack -> sack.complete().contains(item)))
                .max((a, b) -> priority.get(a).compareTo(priority.get(b)));

        return as.stream();
    }

    private record Rucksack(String left, String right, String complete) {
    }
}
