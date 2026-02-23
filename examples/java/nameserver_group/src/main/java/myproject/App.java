package myproject;

import com.pulumi.Pulumi;
import io.github.kitstream.netbird.*;
import io.github.kitstream.netbird.inputs.*;
import java.util.List;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var group = new Group("example-group", GroupArgs.builder()
                .name("Example Java nameserver_group Group")
                .build());

            var res = new NameserverGroup("test-nameserver-group", NameserverGroupArgs.builder()
                .name("Pulumi Java NS Group")
                .description("Pulumi Java NS Group")
                .enabled(true)
                .nameservers(List.of(NameserverGroupNameserverArgs.builder().ip("1.1.1.1").port(53).build()))
                .groups(group.id().applyValue(List::of))
                .primary(true)
                .build());

            ctx.export("resourceName", res.name());
        });
    }
}
