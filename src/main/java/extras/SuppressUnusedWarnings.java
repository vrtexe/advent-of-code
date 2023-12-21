package extras;

import java.beans.Transient;
import java.lang.annotation.*;

import static java.lang.annotation.ElementType.*;

@Documented
@Target({TYPE, FIELD, METHOD, PARAMETER, CONSTRUCTOR, LOCAL_VARIABLE, MODULE})
@Inherited
@Retention(RetentionPolicy.SOURCE)
@SuppressWarnings("unused")
public @interface SuppressUnusedWarnings {
}
