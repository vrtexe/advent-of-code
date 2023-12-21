package solutions.year2022;

import common.Problem;

import java.util.Arrays;
import java.util.Collections;
import java.util.Map;
import java.util.concurrent.atomic.AtomicLong;
import java.util.stream.Stream;

public class Day25 extends Problem {

    public static final String file = "2022/25.txt";
    public static final Problem INSTANCE = new Day25(file, SolutionMethod.PART_ONE);

    public final Long SNAFU_BASE = 5L;
    public final Map<String, Long> snafuDigits = Map.of(
            "=", -2L,
            "-", -1L,
            "0", 0L,
            "1", 1L,
            "2", 2L);

    public Day25(String file) {
        super(file);
    }

    public Day25(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    @Override
    public String solutionPartOne(Stream<String> data) {

        return decimal2Snafu(data.mapToLong(this::snafu2Decimal).sum());
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        // TODO Auto-generated method stub
        return "";
    }

    public Long snafu2Decimal(String snafu) {
        final var counter = new AtomicLong();
        final var snafuLetters = Arrays.asList(snafu.split(""));

        Collections.reverse(snafuLetters);

        return snafuLetters.stream()
                .mapToLong(digit -> snafuDigits.get(digit))
                .map(digit -> digit * (long) Math.pow((double) SNAFU_BASE.longValue(), (double) counter.getAndIncrement()))
                .sum();
    }

    public String decimal2Snafu(Long decimal) {

        var dec = decimal;
        var result = new StringBuilder();

        while (dec != 0) {
            final var div = dec / SNAFU_BASE;
            final var mod = dec % SNAFU_BASE;

            if (mod <= 2) {
                result.insert(0, mod);
                dec = div;
                continue;
            }

            if (mod == 3) {
                result.insert(0, "=");
                dec = div + 1;
                continue;
            }

            if (mod == 4) {
                result.insert(0, "-");
                dec = div + 1;
                continue;
            }
        }

        return result.toString();
    }

}
