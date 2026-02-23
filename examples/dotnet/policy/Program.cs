using System.Collections.Generic;
using Pulumi;
using KitStream.Pulumi.Netbird;

return await Deployment.RunAsync(() => 
{
    var group = new Group("example-group", new GroupArgs
    {
        Name = "Example DotNet policy Group",
    });

    var res = new Policy("test-policy", new PolicyArgs
    {
        Name = "Pulumi DotNet Policy",
        Enabled = true,
        Rule = new PolicyRuleArgs { Action = "accept", Enabled = true, Name = "rule1", Sources = { group.Id }, Destinations = { group.Id } },
    });

    return new Dictionary<string, object?>
    {
        ["resourceName"] = res.Name,
    };
});
