using System.Collections.Generic;
using Pulumi;
using KitStream.Pulumi.Netbird;

return await Deployment.RunAsync(() => 
{
    var group = new Group("example-group", new GroupArgs
    {
        Name = "Example DotNet setup_key Group",
    });

    var res = new SetupKey("test-setup-key", new SetupKeyArgs
    {
        Name = "Pulumi DotNet Setup Key",
        Type = "reusable",
        ExpirySeconds = 86400,
        AutoGroups = { group.Id },
    });

    return new Dictionary<string, object?>
    {
        ["resourceName"] = res.Name,
    };
});
