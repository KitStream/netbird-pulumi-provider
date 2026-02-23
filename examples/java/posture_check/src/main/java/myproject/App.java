package myproject;

import com.pulumi.Pulumi;
import io.github.kitstream.netbird.*;
import io.github.kitstream.netbird.inputs.*;
import java.util.List;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var group = new Group("example-group", GroupArgs.builder()
                .name("Example Java posture_check Group")
                .build());

            var res = new PostureCheck("test-posture-check", PostureCheckArgs.builder()
                .name("Pulumi Java Posture Check")
                .osVersionCheck(PostureCheckOsVersionCheckArgs.builder().darwinMinVersion("1.0.0").build())
                .build());

            ctx.export("resourceName", res.name());
        });
    }
}
