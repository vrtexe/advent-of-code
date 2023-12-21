package solutions.year2022;

import common.Problem;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.math.RoundingMode;
import java.util.*;
import java.util.function.Function;
import java.util.regex.Pattern;
import java.util.stream.Collector;
import java.util.stream.Collectors;
import java.util.stream.LongStream;
import java.util.stream.Stream;

@SuppressWarnings("unused")
public class Day11 extends Problem {

    public static final String file = "2022/11.txt";
    public static final Problem INSTANCE = new Day11(file);

    private static final boolean visual = false;

    public Day11(String file) {
        super(file);
    }

    public Day11(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    private static List<String> split(String line, String delim) {
        return Arrays.stream(line.split(delim)).map(String::trim).toList();
    }

    private static <T> Collector<T, ?, Stack<T>> toReverseStack() {
        return Collectors.collectingAndThen(
                Collectors.toCollection(Stack::new), col -> {
                    Collections.reverse(col);
                    return col;
                });
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    @Override
    public String solutionPartOne(Stream<String> data) {
        final var parser = new MonkeyParser();
        final var game = parser.parseData(data.collect(Collectors.joining("\n")));

        game.play();

        return String.valueOf(game.getMonkeyBusiness());
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        final var parser = new MonkeyParser();
        final var game = parser.parseData(data.collect(Collectors.joining("\n")), 10000L, 1L);

        game.play();

        return String.valueOf(game.getMonkeyBusiness());
    }

    private enum MonkeyAttribute {
        OPERATION("Operation"),
        TEST("Test"),
        TRUE("If true"),
        FALSE("If false"),
        ITEMS("Starting items");

        private static final Map<String, MonkeyAttribute> monkeyAttributeMap = Arrays.stream(values())
                .collect(Collectors.toMap(MonkeyAttribute::getValue, instruction -> instruction));
        private final String value;

        MonkeyAttribute(String value) {
            this.value = value;
        }

        public static MonkeyAttribute fromValue(String value) {
            return monkeyAttributeMap.get(value);
        }

        public String getValue() {
            return value;
        }
    }

    public enum Operation {
        DIVIDE("/"),
        MULTIPLY("*"),
        ADD("+"),
        SUBTRACT("-");

        public final String value;

        Operation(String value) {
            this.value = value;
        }

        public static List<Operation> getValues() {
            return List.of(Operation.values());
        }

        public static Operation find(String expression) {
            return Operation.getValues().stream()
                    .filter(delim -> expression.split(delim.getEscapedValue()).length > 1).findFirst()
                    .orElse(null);
        }

        public String getValue() {
            return this.value;
        }

        public String getEscapedValue() {
            return "\\" + this.getValue();
        }
    }

    private record Monkey(Long number,
                          Function<BigInteger, BigInteger> operation,
                          Stack<Item> items,
                          Function<BigInteger, Boolean> test,
                          Map<Boolean, Long> throwCases) {

        private static final MonkeyPrinter printer = new MonkeyPrinter(visual);

        public BigInteger divideAndRoundDown(BigInteger worry, Long divisor) {
            return new BigDecimal(worry)
                    .divide(new BigDecimal(divisor), RoundingMode.DOWN)
                    .toBigInteger();
        }

        public Item inspectItem(Item item, Long divisor, BigInteger commonDivisor) {
            final var newWorry = operation.apply(item.worryLevel);

            printer.printComputedWorry(newWorry);

            final var boredWorryLevel = divideAndRoundDown(newWorry, divisor).mod(commonDivisor);

            printer.printBoredLevels(divisor, boredWorryLevel);

            return new Item(boredWorryLevel);
        }

    }

    private record Item(BigInteger worryLevel) {
    }

    private static class MonkeyPrinter {
        public final boolean visual;

        public MonkeyPrinter(Boolean visual) {
            this.visual = visual;
        }

        public void printMonkeyTitle(Monkey monkey) {
            if (!visual) {
                return;
            }

            System.out.printf("Monkey %s:%n", monkey.number());
        }

        public void printPossessions(Collection<Monkey> monkeys) {
            if (!visual) {
                return;
            }

            System.out.println(monkeyPossessionsToString(monkeys));
        }

        public String monkeyPossessionsToString(Collection<Monkey> monkeys) {
            return monkeys.stream()
                    .map(monkey -> "Monkey %s: %s"
                            .formatted(monkey.number(), itemsToString(monkey.items())))
                    .collect(Collectors.joining("\n"));
        }

        public String itemsToString(Collection<Item> items) {
            return items.stream()
                    .map(item -> item.worryLevel().toString())
                    .collect(Collectors.joining(", "));
        }

        public void printInspections(Map<Long, Long> inspections) {
            if (!visual) {
                return;
            }

            System.out.println(inspectionsToString(inspections));
        }

        public String inspectionsToString(Map<Long, Long> inspections) {
            return inspections.entrySet()
                    .stream()
                    .map((entry) -> "Monkey %s inspected items %s times."
                            .formatted(entry.getKey(), entry.getValue()))
                    .collect(Collectors.joining("\n"));
        }

        public void printItemInspection(Item item) {
            if (!visual) {
                return;
            }
            System.out.println("Monkey inspects an item with a worry level of %s.".formatted(item.worryLevel).indent(2));
        }

        public void printTestCaseResult(Monkey monkey, Item item) {
            if (!visual) {
                return;
            }
            System.out.println("Current worry level is%s divisible by monkey divisor."
                    .formatted(monkey.test().apply(item.worryLevel()) ? "" : " not")
                    .indent(4));
        }

        public void printThrown(Long monkeyId, Item item) {
            if (!visual) {
                return;
            }

            System.out.println("Item with worry level %s is thrown to monkey %s."
                    .formatted(item.worryLevel(), monkeyId)
                    .indent(4));
        }

        private void printComputedWorry(BigInteger worry) {
            if (!visual) {
                return;
            }
            System.out.println("Worry level is calculated to be %s.".formatted(worry).indent(4));
        }

        private void printBoredLevels(Long divisor, BigInteger boredWorryLevel) {
            if (!visual) {
                return;
            }
            System.out.println(
                    "Monkey gets bored with item. Worry level is divided by %s to %s.".formatted(divisor, boredWorryLevel)
                            .indent(4));
        }
    }

    private static class KeepAway {
        private final MonkeyPrinter monkeyPrinter = new MonkeyPrinter(visual);

        private final long rounds;
        private final long divisor;
        private final BigInteger commonDivisor;

        private Map<Long, Monkey> monkeysMap;
        private Queue<Monkey> monkeys;

        private Map<Long, Long> inspections;

        public KeepAway(Queue<Monkey> monkeys, Long rounds, Long divisor, BigInteger commonDivisor) {
            this.initCommon(monkeys);

            this.rounds = rounds;
            this.divisor = divisor;
            this.commonDivisor = commonDivisor;
        }

        public KeepAway(Queue<Monkey> monkeys, BigInteger commonDivisor) {
            this.initCommon(monkeys);

            this.rounds = 20L;
            this.divisor = 3L;
            this.commonDivisor = commonDivisor;
        }

        private void initCommon(Queue<Monkey> monkeys) {
            this.monkeys = monkeys;

            this.monkeysMap = monkeys.stream()
                    .collect(Collectors.toMap(Monkey::number, monkey -> monkey));

            this.inspections = monkeysMap.keySet().stream()
                    .collect(Collectors.toMap(monkeyNumber -> monkeyNumber, __ -> 0L));
        }

        public void play() {
            LongStream.range(0, rounds).forEach(_num -> {
                playOutRound();
                monkeyPrinter.printPossessions(monkeys);
            });
            monkeyPrinter.printInspections(inspections);
        }

        public Long getMonkeyBusiness() {
            return inspections.values()
                    .stream().sorted(Comparator.reverseOrder())
                    .limit(2)
                    .reduce(1L, (a, b) -> a * b);
        }

        private void playOutRound() {
            monkeys.forEach(monkey -> {
                monkeyPrinter.printMonkeyTitle(monkey);
                inspectAllMonkeyItems(monkey);
            });
        }

        private void inspectAllMonkeyItems(Monkey monkey) {
            while (!monkey.items().empty()) {
                inspectItem(monkey, monkey.items.pop());
            }
        }

        private void inspectItem(Monkey monkey, Item item) {
            monkeyPrinter.printItemInspection(item);

            final var newItem = monkey.inspectItem(item, divisor, commonDivisor);
            final var newMonkey = monkey.throwCases()
                    .get(monkey.test().apply(newItem.worryLevel()));

            monkeyPrinter.printTestCaseResult(monkey, newItem);

            throwTo(newMonkey, newItem);
            recordInspection(monkey);
        }

        private void throwTo(Long monkeyId, Item item) {
            monkeyPrinter.printThrown(monkeyId, item);
            monkeysMap.get(monkeyId).items.push(item);
        }

        private void recordInspection(Monkey monkey) {
            inspections.computeIfPresent(monkey.number(), (_monkey, inspections) -> inspections + 1L);
        }
    }

    @SuppressWarnings("UnusedReturnValue")
    private static class MonkeyParser {

        private final List<BigInteger> divisors = new ArrayList<>();

        public KeepAway parseData(String data) {
            final var monkeys = parseMonkeys(data);

            return new KeepAway(monkeys, this.getCommonDivisor());
        }

        public KeepAway parseData(String data, Long rounds, Long divisor) {
            final var monkeys = parseMonkeys(data);

            return new KeepAway(monkeys, rounds, divisor, this.getCommonDivisor());
        }

        public BigInteger getCommonDivisor() {
            return divisors.stream()
                    .reduce(new BigInteger("1"), BigInteger::multiply);
        }

        private void saveDivisor(BigInteger divisor) {
            this.divisors.add(divisor);
        }

        private Queue<Monkey> parseMonkeys(String data) {
            return split(data, "\\n\\n")
                    .stream()
                    .map(this::parseMonkey)
                    .collect(Collectors.toCollection(LinkedList::new));
        }

        private Monkey parseMonkey(String monkey) {
            final var monkeyAttributeBuilder = new MonkeyAttributeBuilder();

            Arrays.stream(monkey.split("\n"))
                    .map(String::trim)
                    .peek(line -> extractMonkeyNumber(line, monkeyAttributeBuilder))
                    .skip(1)
                    .forEach(line -> extractMonkeyAttributes(line, monkeyAttributeBuilder));

            return monkeyAttributeBuilder.createMonkey();
        }

        private MonkeyAttributeBuilder extractMonkeyAttributes(String line, MonkeyAttributeBuilder monkeyAttributes) {
            final var attribute = split(line, ":");

            switch (MonkeyAttribute.fromValue(attribute.get(0))) {
                case OPERATION -> monkeyAttributes.operation(parseOperation(attribute.get(1)));
                case TEST -> monkeyAttributes.test(parseTestCase(attribute.get(1)));
                case ITEMS -> monkeyAttributes.items(parseItems(attribute.get(1)));
                case TRUE -> monkeyAttributes.throwCase(parseThrowCase(attribute.get(1), true));
                case FALSE -> monkeyAttributes.throwCase(parseThrowCase(attribute.get(1), false));
            }

            return monkeyAttributes;
        }

        private void extractMonkeyNumber(String line, MonkeyAttributeBuilder monkeyAttributeBuilder) {
            final var pattern = Pattern.compile("Monkey ([0-9]+):");
            final var matcher = pattern.matcher(line);

            if (matcher.matches()) {
                final var monkeyId = Long.parseLong(matcher.group(1));
                monkeyAttributeBuilder.monkeyId(monkeyId);
            }
        }

        private Stack<Item> parseItems(String items) {
            return split(items, ",").stream()
                    .map(BigInteger::new)
                    .map(Item::new)
                    .collect(toReverseStack());
        }

        private Function<BigInteger, Boolean> parseTestCase(String testCase) {
            final var pattern = Pattern.compile("divisible by ([0-9]+)");
            final var matcher = pattern.matcher(testCase);

            if (matcher.find()) {
                final var divisor = new BigInteger(matcher.group(1));

                saveDivisor(divisor);

                return (number) -> number.mod(divisor).equals(new BigInteger("0"));
            }

            return null;
        }

        private MonkeyAttributeBuilder.ThrowCase parseThrowCase(String monkey, Boolean choice) {
            final var pattern = Pattern.compile("throw to monkey ([0-9]+)");
            final var matcher = pattern.matcher(monkey);

            if (!matcher.find()) {
                throw new RuntimeException("Throw expression does not match");
            }

            final var monkeyId = Long.parseLong(matcher.group(1));

            return new MonkeyAttributeBuilder.ThrowCase(choice, monkeyId);
        }

        private Function<BigInteger, BigInteger> parseOperation(String operation) {
            final var sign = Operation.find(operation);

            final var sides = split(operation, "=");
            final var args = split(sides.get(1), sign.getEscapedValue());

            return (old) -> evaluateExpression(getArgOrDefault(args.get(0), old), sign, getArgOrDefault(args.get(1), old));
        }

        private BigInteger evaluateExpression(BigInteger left, Operation operation, BigInteger right) {
            return switch (operation) {
                case ADD -> left.add(right);
                case SUBTRACT -> left.subtract(right);
                case MULTIPLY -> left.multiply(right);
                case DIVIDE -> left.divide(right);
            };
        }

        @SuppressWarnings("SwitchStatementWithTooFewBranches")
        private BigInteger getArgOrDefault(String arg, BigInteger old) {
            return switch (arg) {
                case "old" -> old;
                default -> new BigInteger(arg);
            };
        }
    }

    @SuppressWarnings("UnusedReturnValue")
    private static class MonkeyAttributeBuilder {
        Long monkeyId;
        Function<BigInteger, BigInteger> operation;
        Function<BigInteger, Boolean> test;
        Map<Boolean, Long> throwCase = new HashMap<>();
        Stack<Item> items = new Stack<>();

        public Monkey createMonkey() {
            return new Monkey(monkeyId, operation, items, test, throwCase);
        }

        public MonkeyAttributeBuilder monkeyId(Long monkeyId) {
            this.monkeyId = monkeyId;
            return this;
        }

        public MonkeyAttributeBuilder operation(Function<BigInteger, BigInteger> operation) {
            this.operation = operation;
            return this;
        }

        public MonkeyAttributeBuilder test(Function<BigInteger, Boolean> test) {
            this.test = test;
            return this;
        }

        public MonkeyAttributeBuilder throwCase(ThrowCase throwCase) {
            this.throwCase.put(throwCase.condition(), throwCase.monkeyId());
            return this;
        }

        public MonkeyAttributeBuilder items(Stack<Item> items) {
            this.items = items;
            return this;
        }

        private record ThrowCase(Boolean condition, Long monkeyId) {
        }
    }
}
