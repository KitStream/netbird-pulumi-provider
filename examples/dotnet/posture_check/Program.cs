using System.Collections.Generic;
using Pulumi;
using KitStream.Pulumi.Netbird;

return await Deployment.RunAsync(() => 
{
    var group = new Group("example-group", new GroupArgs
    {
        Name = "Example DotNet posture_check Group",
    });

    var res = new PostureCheck("test-posture-check", new PostureCheckArgs
    {
        Name = "Pulumi DotNet Posture Check",
        OsVersionCheck = new PostureCheckOsVersionCheckArgs { DarwinMinVersion = "1.0.0" },
    });

    return new Dictionary<string, object?>
    {
        ["resourceName"] = res.Name,
    };
});
