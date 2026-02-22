package com.minimal;

import com.pulumi.Pulumi;
import io.github.kitstream.netbird.Group;
import io.github.kitstream.netbird.GroupArgs;
import java.util.Map;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var group = new Group("test-group", GroupArgs.builder()
                .name("Pulumi Java Test Group")
                .build());

            ctx.export("groupName", group.name());
        });
    }
}
