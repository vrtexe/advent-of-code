package solutions.year2022;

import common.Problem;

import java.util.*;
import java.util.function.Consumer;
import java.util.stream.Collectors;
import java.util.stream.LongStream;
import java.util.stream.Stream;

@SuppressWarnings("unused")
public class Day10 extends Problem {

    public static final String file = "2022/10.txt";
    public static final Problem INSTANCE = new Day10(file, SolutionMethod.ALL);

    public Day10(String file) {
        super(file);
    }

    public Day10(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    @Override
    public String solutionPartOne(Stream<String> data) {
        final var commandRunner = new CommandRunner();
        data.map(line -> Arrays.asList(line.split(" ")))
                .forEach(commandRunner::addCommand);

        commandRunner.executeCommands();

        return String.valueOf(commandRunner.getResult());
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        final var commandRunner = new CommandRunner();
        data.map(line -> Arrays.asList(line.split(" ")))
                .forEach(commandRunner::addCommand);

        commandRunner.executeCommands();

        return commandRunner.getScreenState();
    }

    private enum Instruction {
        NOOP("noop"),
        ADDX("addx");

        private static final Map<String, Instruction> InstructionValueMapping = Arrays.stream(values())
                .collect(Collectors.toMap(Instruction::getValue, instruction -> instruction));
        private final String value;

        Instruction(String value) {
            this.value = value;
        }

        public static Instruction fromValue(String value) {
            return InstructionValueMapping.get(value);
        }

        public String getValue() {
            return value;
        }
    }

    private interface Command {
        void execute(CommandRunner commandRunner, Consumer<Long> cycle);
    }

    private record NoopCommand(Instruction instruction, List<String> arguments) implements Command {

        private static final int cycles = 1;

        public static NoopCommand of(List<String> instruction) {
            final var args = instruction.size() > 1
                    ? instruction.subList(1, instruction.size())
                    : null;

            return new NoopCommand(Instruction.fromValue(instruction.get(0)), args);
        }

        @Override
        public void execute(CommandRunner commandRunner, Consumer<Long> cycle) {
            cycle.accept(commandRunner.x);
        }
    }

    private record AddxCommand(Instruction instruction, List<Long> arguments) implements Command {

        private static final int cycles = 2;

        public static AddxCommand of(List<String> instruction) {

            return new AddxCommand(Instruction.fromValue(instruction.get(0)), List.of(Long.parseLong(instruction.get(1))));
        }

        @Override
        public void execute(CommandRunner commandRunner, Consumer<Long> cycle) {
            for (int i = 0; i < cycles - 1; i++) {
                cycle.accept(commandRunner.x);
            }

            cycle.accept(commandRunner.x + arguments.get(0));
        }
    }

    private static class Screen {
        private final long rowCycles = 40L;

        private final String off = ".";
        private final String on = "#";

        private final Map<Long, Map<Long, String>> screen;

        private Long spritePosition = 1L;

        public Screen() {
            long rows = 6;
            this.screen = LongStream.range(0, rows).boxed()
                    .collect(Collectors.toMap(row -> row, row -> initRow()));
        }

        private Map<Long, String> initRow() {
            return new TreeMap<>(
                    LongStream.range(0, rowCycles).boxed()
                            .collect(Collectors.toMap(col -> col, col -> off)));
        }

        public void drawPixel(Long cycle) {
            if (isSpriteOnTopOfPixel(getRow(cycle))) {
                screen.computeIfPresent(getPixelPosition(cycle), (k, v) -> {
                    v.computeIfPresent(getRow(cycle), (k1, v1) -> on);
                    return v;
                });
            }
        }

        @SuppressWarnings("all")
        private Long getPixelPosition(Long pixel) {
            return Math.round(Math.ceil(pixel / rowCycles));
        }

        private Long getRow(Long pixel) {
            if (pixel == 0L) {
                return 0L;
            }

            return pixel % rowCycles;
        }

        public Boolean isSpriteOnTopOfPixel(Long pixel) {
            return pixel >= spritePosition - 1 && pixel <= spritePosition + 1;
        }

        public void updateSpritePosition(Long position) {
            this.spritePosition = position;
        }

        @Override
        public String toString() {
            final var stringBuilder = new StringBuilder("\n");
            screen.values().forEach(s -> {
                s.values().forEach(stringBuilder::append);
                stringBuilder.append("\n");
            });

            return stringBuilder.toString();
        }
    }

    private static class CommandRunner {
        private static final long saveValue = 20L;
        private static final long everyValue = 40L;
        private final List<Command> commands = new ArrayList<>();
        private final List<Long> signalStrengths = new ArrayList<>();
        private final Screen screen = new Screen();
        private Long x = 1L;
        private Long cycle = 0L;

        public void addCommand(List<String> instruction) {
            final var command = switch (Instruction.fromValue(instruction.get(0))) {
                case NOOP -> NoopCommand.of(instruction);
                case ADDX -> AddxCommand.of(instruction);
            };

            commands.add(command);
        }

        public void executeCommands() {
            for (var command : commands) {
                command.execute(this, this::cycle);
            }
            // commands.forEach(command -> );
        }

        public Long getResult() {
            return signalStrengths.stream().mapToLong(l -> l).sum();
        }

        public String getScreenState() {
            return this.screen.toString();
        }

        private void cycle(Long value) {
            screen.drawPixel(cycle);

            this.cycle++;

            if (cycle % everyValue == saveValue) {
                signalStrengths.add(x * cycle);
            }

            x = value;
            screen.updateSpritePosition(x);
        }
    }
}
