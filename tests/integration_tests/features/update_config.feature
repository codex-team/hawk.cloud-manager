Feature: Update WireGuard config

  Scenario: Update config
    When I send PUT request to "/config" with data:
    """
    {
      "hosts": [
        {
          "name": "agent1",
          "public_key": "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=",
          "endpoint": "10.5.0.5:12345",
          "allowed_ips": ["10.5.0.0/16"]
        },
        {
          "name": "agent2",
          "public_key": "cXdlcnFld3RyZ3dydHd2cnd0cnRuYnJlcXdlcnFxd3Q=",
          "endpoint": "10.5.0.6:12346",
          "allowed_ips": ["10.5.0.0/16", "234.122.178.0/32"]
        }
      ],
      "groups": [
        {
          "name": "test-group",
          "hosts": ["agent1", "agent2"]
        }
      ]
    }
    """
    Then The response code should be 200
    When I send POST request to "/topology" with data:
    """
    {
      "public_key": "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk="
    }
    """
    Then The response code should be 200
    And I receive data:
    """
    {
      "listen_port": 12345,
      "peers": [
        {
          "public_key": "cXdlcnFld3RyZ3dydHd2cnd0cnRuYnJlcXdlcnFxd3Q=",
          "endpoint": "10.5.0.6:12346",
          "allowed_ips": ["10.5.0.0/16", "234.122.178.0/32"]
        }
      ]
    }
    """
