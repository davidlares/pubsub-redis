class Product < ApplicationRecord
  # after_create represents the method after the object is stored
  after_create :pubsub
  # calling redis client for publish on a selected channel
  def pubsub
    msg = {
      id: self.id,
      name: self.name
    }
    # $redis.publish 'test1', 'Hello,Rails'
    $redis.publish 'test1', msg.to_json
  end
end
