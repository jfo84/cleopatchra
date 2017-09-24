require_relative './data_mixin'

class Review
  include Virtus.model
  include DataMixin

  attribute :pull_id, String
  
  def record
    connection.exec("INSERT INTO reviews (id, data, pull_id)
                    VALUES (#{id}, E'#{data}', #{pull_id})")
  end
end