package solutions.year2022;

import com.google.gson.*;
import common.Problem;

import java.util.*;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.stream.Collectors;
import java.util.stream.Stream;

@SuppressWarnings("unused")
public class Day13 extends Problem {

    public static final String file = "2022/13.txt";
    public static final Problem INSTANCE = new Day13(file, SolutionMethod.ALL);

    public static final Boolean visual = false;

    private final Gson json = new GsonBuilder()
            .registerTypeAdapter(Packet.class, new PacketDeserializer())
            .create();

    public Day13(String file) {
        super(file);
    }

    public Day13(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    @Override
    public String solutionPartOne(Stream<String> data) {
        final var packetPairs = parsePacketPairs(data);

        final var index = new AtomicInteger(0);

        final var indicesSum = packetPairs.stream()
                .map(s -> s.getRightOrderPacketsIndex(index.incrementAndGet())).filter(Objects::nonNull)
                .mapToLong(s -> s)
                .sum();

        return String.valueOf(indicesSum);
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        final var dividerPackets = Stream.of("[[2]]", "[[6]]")
                .map(this::deserializePacket)
                .toList();
        final var packetPairs = Stream.concat(parsePackets(data), dividerPackets.stream());

        final var dividers = packetPairs.sorted(Comparator.reverseOrder()).toList();

        final var divider1 = dividers.indexOf(dividerPackets.get(0)) + 1;
        final var divider2 = dividers.indexOf(dividerPackets.get(1)) + 1;

        return String.valueOf(divider1 * divider2);
    }

    private Stream<Packet> parsePackets(Stream<String> data) {
        return data.filter(s -> !s.isBlank())
                .map(this::deserializePacket);
    }

    private List<PacketPair> parsePacketPairs(Stream<String> data) {
        return Arrays.stream(data.collect(Collectors.joining("\n")).split("\\n\\n"))
                .map(this::parsePacketPair)
                .toList();
    }

    private PacketPair parsePacketPair(String packetPair) {
        List<Packet> packets = Arrays.stream(packetPair.split("\n"))
                .map(this::deserializePacket)
                .toList();

        return new PacketPair(packets.get(0), packets.get(1));
    }

    private Packet deserializePacket(String list) {
        return json.fromJson(list, Packet.class);
    }

    private enum Type {
        LIST,
        INTEGER,
        EMPTY
    }

    private record PacketPair(Packet left, Packet right) {
        private final static PacketComparisonPrinter packetPairPrinter = PacketComparisonPrinter.getInstance();

        private Integer getRightOrderPacketsIndex(Integer index) {
            packetPairPrinter.printPairTitle(index);

            if (this.isOrderGood()) {
                return index;
            }
            return null;
        }

        public Boolean isOrderGood() {
            return left.compare(right);
        }
    }

    private static class PacketComparisonPrinter {
        private final static PacketComparisonPrinter INSTANCE = new PacketComparisonPrinter(visual);

        private final Boolean isVisual;

        private PacketComparisonPrinter(Boolean visual) {
            this.isVisual = visual;
        }

        public static PacketComparisonPrinter getInstance() {
            return INSTANCE;
        }

        public void printPairTitle(Integer index) {
            if (!isVisual) {
                return;
            }

            System.out.println();
            System.out.printf("== Pair %s ==%n", index);
        }

        public void printIntCompare(Packet left, Packet right, Integer level) {
            if (!isVisual) {
                return;
            }

            System.out.printf("- Compare %s vs %s".indent(level), left.number, right.number);
        }

        public void printArrayCompare(Packet left, Packet right, Integer level) {
            if (!isVisual) {
                return;
            }

            System.out.printf("- Compare %s vs %s".indent(level), left.numbers, right.numbers);
        }

        public void printIntCompareResult(Boolean result, Integer level) {
            if (!isVisual) {
                return;
            }

            if (result != null) {
                if (result) {
                    System.out.print("- Left side is smaller, so inputs are in the right order".indent(level));
                } else {
                    System.out.print("- Right side is smaller, so inputs are not in the right order".indent(level));
                }
            }
        }

        public void printRightRanOut(Integer level) {
            if (!isVisual) {
                return;
            }

            System.out.print("- Right side ran out of items, so inputs are not in the right order".indent(level));
        }

        public void printLeftRanOut(Integer level) {
            if (!isVisual) {
                return;
            }

            System.out.print("- Left side ran out of items, so inputs are in the right order".indent(level + 1));
        }

        public void printLeftListRightInt(Packet left, Packet right, Integer level) {
            if (!isVisual) {
                return;
            }

            System.out.printf("- Compare %s vs %s".indent(level), left.numbers, right.number);
        }

        public void printRightListLeftInt(Packet left, Packet right, Integer level) {
            if (!isVisual) {
                return;
            }

            System.out.printf("- Compare %s vs %s".indent(level), left.number, right.numbers);
        }
    }

    private static class Packet implements Comparable<Packet> {
        private final PacketComparisonPrinter packetComparisonPrinter = PacketComparisonPrinter.getInstance();

        private Integer number;
        private List<Packet> numbers = new ArrayList<>();
        private Type type;

        public Packet() {
        }

        public Packet(Integer number, List<Packet> numbers, Type type) {
            this.number = number;
            this.numbers = numbers;
            this.type = type;
        }

        public Packet number(Integer number) {
            this.number = number;
            return this;
        }

        public Packet numbers(Packet numbers) {
            this.numbers.add(numbers);
            return this;
        }

        public Packet type(Type type) {
            this.type = type;
            return this;
        }

        public Boolean compare(Packet that) {
            return getComparisonResult(this.compare(that, 0));
        }

        public int compare(Packet that, Integer level) {
            if (this.type == Type.INTEGER && that.type == Type.INTEGER) {
                packetComparisonPrinter.printIntCompare(this, that, level);

                final var result = compareNumbers(that);
                packetComparisonPrinter.printIntCompareResult(result, level + 2);

                return getComparisonResult(result);
            }

            if (this.type == Type.LIST && that.type == Type.LIST) {
                packetComparisonPrinter.printArrayCompare(this, that, level);
                final var size = Math.max(this.numbers.size(), that.numbers.size());
                for (int i = 0; i < size; i++) {
                    final var thisPacket = this.tryGetIndex(i);
                    final var thatPacket = that.tryGetIndex(i);

                    if (thisPacket == null && thatPacket == null) {
                        continue;
                    }

                    if (thatPacket == null) {
                        packetComparisonPrinter.printRightRanOut(level + 2);
                        return getComparisonResult(false);
                    }

                    if (thisPacket == null) {
                        packetComparisonPrinter.printLeftRanOut(level + 2);
                        return getComparisonResult(true);
                    }

                    final var result = thisPacket.compare(thatPacket, level + 2);
                    if (result != getComparisonResult(null)) {
                        return result;
                    }
                }

                return getComparisonResult(null);
            }

            if (this.type == Type.LIST) {
                packetComparisonPrinter.printLeftListRightInt(this, that, level);
                return this.compare(new Packet(null, List.of(that), Type.LIST), level + 2);
            }

            if (that.type == Type.LIST) {
                packetComparisonPrinter.printRightListLeftInt(this, that, level);
                return new Packet(null, List.of(this), Type.LIST).compare(that, level + 2);
            }

            return getComparisonResult(true);
        }

        private int getComparisonResult(Boolean value) {
            if (value == null) {
                return 0;
            }

            if (value) {
                return 1;
            }

            return -1;
        }

        private Boolean getComparisonResult(int value) {
            return switch (value) {
                case 0, 1 -> true;
                case -1 -> false;
                default -> null;
            };
        }

        private Packet tryGetIndex(Integer index) {
            try {
                return this.numbers.get(index);
            } catch (Exception e) {
                return null;
            }
        }

        private Boolean compareNumbers(Packet other) {
            if (this.number.equals(other.number)) {
                return null;
            }

            return this.number < other.number;

        }

        @Override
        public String toString() {
            if (type == Type.INTEGER) {

                return "%s".formatted(number != null ? number : "");
            }

            if (type == Type.LIST) {
                return "%s".formatted(numbers != null ? numbers : "");
            }

            return "[]";
        }

        @Override
        public int compareTo(Packet that) {
            return this.compare(that, 0);
        }

    }

    private static class PacketDeserializer implements JsonDeserializer<Packet> {

        @Override
        public Packet deserialize(JsonElement json, java.lang.reflect.Type typeOfT, JsonDeserializationContext context)
                throws JsonParseException {

            final var builder = new Packet();

            record BuilderJsonBind(Packet builder, JsonElement element) {
            }

            final var builderQueue = new LinkedList<BuilderJsonBind>();
            builderQueue.offer(new BuilderJsonBind(builder, json));

            while (!builderQueue.isEmpty()) {
                final var currentBuilder = builderQueue.poll();

                if (!currentBuilder.element().isJsonArray()) {
                    final var number = tryGetNumber(currentBuilder.element());
                    currentBuilder.builder().number(number).type(Type.INTEGER);
                } else {
                    if (currentBuilder.element().getAsJsonArray().isEmpty()) {
                        currentBuilder.builder().number(null)
                                .numbers(null).type(Type.LIST);
                    }
                    for (final var item : currentBuilder.element().getAsJsonArray()) {
                        final var nestedBuilder = new Packet();
                        builderQueue.offer(new BuilderJsonBind(nestedBuilder, item));
                        currentBuilder.builder().numbers(nestedBuilder).type(Type.LIST);
                    }
                }
            }

            return builder;
        }

        private Integer tryGetNumber(JsonElement json) {
            try {
                return json.getAsInt();
            } catch (Exception e) {
                return null;
            }
        }

    }
}
